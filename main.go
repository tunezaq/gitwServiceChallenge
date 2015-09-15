package main

import (
	`bytes`
	`encoding/gob` // Persistence
	`encoding/json`
	`fmt`
	`io/ioutil`
	`log`
	`net/http`
	`os`
	`os/signal`
	`strconv`
	`sync`
	`sync/atomic`
	`syscall`
	`time`
)

var (
	x_cache map[interface{}]interface{}
	x_count map[interface{}]int
	lock    *sync.Mutex
	prefix  string = `/cache/`
	plen    int    = len(prefix)
	wg      *sync.WaitGroup
	stop    chan os.Signal
	sched   int32
)

type (
	cache    map[string]interface{}
	cacheElt struct {
		Key   interface{} `json:"key"`
		Value interface{} `json:"value"`
	}
	flatCache struct {
		Elts []cacheElt `json:"cache"`
	}
)

func scheduleUpdate() {
	atomic.StoreInt32(&sched, 1)
}

func persist() {
	defer wg.Done()
	tick := time.Tick(500 * time.Millisecond)
	p := func() {
		if atomic.CompareAndSwapInt32(&sched, 1, 0) {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			lock.Lock()
			defer lock.Unlock()
			if err := enc.Encode(x_cache); err != nil {
				log.Printf(`[ERROR] Unable to persist cache: %v`, err)
				return
			}
			if err := enc.Encode(x_count); err != nil {
				log.Printf(`[ERROR] Unable to persist counts: %v`, err)
				return
			}
			if err := ioutil.WriteFile(`/tmp/kirkwood.dat`, buf.Bytes(), 0644); err != nil {
				log.Printf(`[ERROR] Unable to write persist file: %v`, err)
			}
			log.Printf(`Persisted %d items.`, len(x_cache))
		}
	}

	for {
		select {
		case <-tick:
			p()
		case <-stop:
			p()
			log.Printf(`Done.`)
			return
		}
	}
}

func unpersist() {
	if b, err := ioutil.ReadFile(`/tmp/kirkwood.dat`); err == nil {
		buf := bytes.NewBuffer(b)
		dec := gob.NewDecoder(buf)
		if err = dec.Decode(&x_cache); err != nil {
			log.Printf("[ERROR] Unable to unpersist: %v\n", err)
			return
		}
		if err = dec.Decode(&x_count); err != nil {
			log.Printf("[ERROR] Unable to unpersist: %v\n", err)
		}
	}
}

func update(k interface{}, v interface{}) int {
	lock.Lock()
	defer lock.Unlock()
	if _, found := x_cache[k]; found {
		x_cache[k] = v
		scheduleUpdate()
		return 204
	}
	return 404
}

func create(k interface{}, v interface{}) int {
	lock.Lock()
	defer lock.Unlock()
	if _, found := x_cache[k]; !found {
		x_cache[k] = v
		x_count[k] = 0
		scheduleUpdate()
		return 201
	}
	return 409
}

func get(k string) ([]cacheElt, int) {
	lock.Lock()
	defer lock.Unlock()
	// get only works on strings (from the URL), so we have to try to
	// coerce other types from the key.
	elts := make([]cacheElt, 0)
	if v, ok := x_cache[k]; ok {
		elts = append(elts, cacheElt{k, v})
		if x_count[k] == 99 {
			delete(x_cache, k)
			delete(x_count, k)
		} else {
			x_count[k]++
		}
	}
	if i, err := strconv.ParseInt(k, 10, 64); err == nil {
		if v, ok := x_cache[int(i)]; ok {
			elts = append(elts, cacheElt{i, v})
			if x_count[i] == 99 {
				delete(x_cache, i)
				delete(x_count, i)
			} else {
				x_count[i]++
			}
		}
	}
	if f, err := strconv.ParseFloat(k, 64); err == nil {
		if v, ok := x_cache[f]; ok {
			elts = append(elts, cacheElt{f, v})
			if x_count[f] == 99 {
				delete(x_cache, f)
				delete(x_count, f)
			} else {
				x_count[f]++
			}
		}
	}
	if b, err := strconv.ParseBool(k); err == nil {
		if v, ok := x_cache[b]; ok {
			elts = append(elts, cacheElt{b, v})
			if x_count[b] == 99 {
				delete(x_cache, b)
				delete(x_count, b)
			} else {
				x_count[b]++
			}
		}
	}

	if len(elts) > 0 {
		scheduleUpdate()
		return elts, 200
	}

	return nil, 404
}

func rm(k string) int {
	lock.Lock()
	defer lock.Unlock()
	// get only works on strings (from the URL), so we have to try to
	// coerce other types from the key.
	ret := 404
	if k == "" {
		x_cache = make(map[interface{}]interface{})
		x_count = make(map[interface{}]int)
		ret = 204
	}
	if _, found := x_cache[k]; found {
		delete(x_cache, k)
		delete(x_count, k)
		ret = 204
	}
	if i, err := strconv.ParseInt(k, 10, 64); err == nil {
		if _, found := x_cache[i]; found {
			delete(x_cache, i)
			delete(x_count, i)
			ret = 204
		}
	}
	if f, err := strconv.ParseFloat(k, 64); err == nil {
		if _, found := x_cache[f]; found {
			delete(x_cache, f)
			delete(x_count, f)
			ret = 204
		}
	}
	if b, err := strconv.ParseBool(k); err == nil {
		if _, found := x_cache[b]; found {
			delete(x_cache, b)
			delete(x_count, b)
			ret = 204
		}
	}

	if ret == 204 {
		scheduleUpdate()
	}

	return ret
}

func flatten() flatCache {
	lock.Lock()
	defer lock.Unlock()
	elts := make([]cacheElt, len(x_cache))
	f := flatCache{elts}
	i := 0
	for k, v := range x_cache {
		f.Elts[i].Key, f.Elts[i].Value = k, v
		i++
	}
	return f
}

func parseArg(b []byte) (elt cacheElt, err error) {
	err = json.Unmarshal(b, &elt)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, bodyerr := ioutil.ReadAll(r.Body)

	var (
		ret   int
		s     interface{}
		key   = r.URL.Path[plen:]
		v     []cacheElt
		abort = false
	)

	switch r.Method {
	case `DELETE`:
		fmt.Printf("TIME TO DELETE\n\n\n")
		ret = rm(key)
	case `GET`:
		if key == `` {
			s = flatten()
			ret = 200
		} else {
			v, ret = get(key)
			if abort = ret == 404; !abort {
				if len(v) == 0 {
					s = nil
				} else if len(v) == 1 {
					s = v
				} else {
					s = flatCache{v}
				}
			}
		}
	case `POST`:
		if bodyerr != nil {
			log.Printf("[ERROR] Reading payload: %v\n", bodyerr)
			ret = 406
			abort = true
		} else if elt, err := parseArg(body); err == nil {
			ret = create(elt.Key, elt.Value)
			abort = ret == 404
			if ret == 201 {
				s = fmt.Sprintf(`/cache/%s`, key)
			}
		}
	case `PUT`:
		if bodyerr != nil {
			fmt.Printf("[ERROR] Reading payload: %v\n", bodyerr)
			ret = 406
			abort = true
		} else if elt, err := parseArg(body); err == nil {
			if fmt.Sprintf(`%v`, elt.Key) != key {
				ret = 406 // Key mismatch (?)
				abort = true
			} else {
				ret = update(elt.Key, elt.Value)
				abort = ret == 404
			}
		}
	}

	body = nil
	if !abort {
		if ret == 201 {
			w.Header().Add(`Location`, s.(string))
		} else {
			body, _ = json.Marshal(s)
		}
	}
	w.WriteHeader(ret)
	if body != nil {
		fmt.Fprintf(w, string(body))
	}
}

func serve() {
	http.HandleFunc(`/cache/`, handler)
	http.ListenAndServe(`:8088`, nil)
}

func main() {
	lock = &sync.Mutex{}
	x_cache = make(map[interface{}]interface{})
	x_count = make(map[interface{}]int)
	stop = make(chan os.Signal, 1)

	unpersist()

	wg = new(sync.WaitGroup)
	wg.Add(1)
	go persist()

	create(`foo`, `bar`)
	create(`baz`, 100000000.000000001)
	create(`quux`, `Hello, world!`)
	create(123, `Integer`)
	create(123.0, `Float`)
	create(false, `Boolean`)
	r, v := get(`123`)
	fmt.Println(`123 =>`, r, v)
	r, v = get(`123.0`)
	fmt.Println(`123.0 =>`, r, v)
	r, v = get(`asdlfkj`)
	fmt.Println(`asdlfkj =>`, r, v)

	go serve()

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-stop
	wg.Wait()
}

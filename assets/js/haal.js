/*haal.post('https://abc.xyz/', {
  body: {
    action: 'update-condition',
    data: '42'
  }
})(res => {
  if (!res.ok) console.error('bad things happened: ', res.err)
  if (res.out) {
    // do something
  } else {
    // check with res for clues as to what went wrong
  }
})*/
{
  const ctype = 'Content-Type'

  var haal = new Proxy((endpoint, options = {}) => {
    if (options == null || typeof options !== 'object') options = {}
    if (endpoint != null && endpoint.constructor === Object) {
      [options, endpoint] = [endpoint, options.endpoint]
      if (!endpoint) throw new Error('road to nowhere? requests need enpoints')
    }
    if (options.method == null) options.method = 'GET'
    if (!options.headers) options.headers = {}
    if (haal.headers) {
      options.headers = Object.assign({}, haal.headers, options.headers)
    }

    if (!options.credential) {
      options.credential = 'include'
    }

    if (typeof options.body === 'object') {
      if (msgpack && !options.json) {
        options.headers[ctype] = 'application/msgpack'
        options.body = msgpack.encode(options.body)
      } else {
        options.headers[ctype] = 'application/json'
        options.body = JSON.stringify(options.body)
      }
    }

    return handle => new Promise((resolve, reject) => {
      fetch(endpoint, options).then(res => {
        const type = res.headers.get(ctype)
        const handleErr = err => {
          res.ok = false
          res.err = err
          if (handle) handle(res, err)
          reject(res)
        }
        const handleOK = out => {
          res.out = out
          if (handle) handle(res)
          resolve(res)
        }
        if (type.includes('/json')) {
          res.json().then(handleOK, handleErr)
        } else if (type.includes('/msgpack')) {
          res.blob().then(blob => {
            try {
              haal.filereader.onerror = null
              haal.filereader.onload = null
              haal.filereader.onerror = handleErr
              haal.filereader.onload = e => {
                try {
                  handleOK(msgpack.decode(new Uint8Array(e.target.result)))
                } catch (err) {
                  handleErr(err)
                } finally {
                  haal.filereader.onload = null
                  haal.filereader.onerror = null
                  haal.filereader.abort()
                }
              }
              haal.filereader.readAsArrayBuffer(blob)
            } catch (e) {
              handleErr(e)
            }
          }, handleErr)
        } else if (type.includes('text/') || type.includes('application/javascript')) {
          res.text().then(handleOK, handleErr)
        }
      })
    })
  }, {
    get(h, key) {
      return key in h ? Reflect.get(h, key) : (endpoint, options = {}) => {
        if (endpoint != null && endpoint.constructor === Object) { 
          [options, endpoint] = [endpoint, options.endpoint]
        }
        options.method = key.toUpperCase()
        return haal(endpoint, options)
      }
    }
  })

  haal.filereader = new FileReader()
  haal.arrayEqual = (a, b) => !(a < b || b < a)
  haal.infinify = (fn, reflect = false) => new Proxy(fn, {
    get: (fn, key) => reflect && key in fn ? Reflect.get(fn, key) : fn.bind(null, key)
  })
}

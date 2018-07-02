# internetzValidator / valid8me

### What
A service that you feed a URL and get back the status code from that URL

Bonus features:
- Simplistic Facebook URL validator
- ~~Simplistic LinkedIn URL validator~~ (Currently mostly broken due to LinkedIn returning a non-standard `999` status code for nearly any request)
- Twitter handle checker
- Instagram username checker

### Why
I have a project requiring basic validation of Facebook and LinkedIn URLs, and Twitter/Instagram usernames.

### How
Check out the live testing instance: [https://valid8me.stag9.com](https://valid8me.stag9.com)

Otherwise, use the provided Dockerfile...
```
docker build -t valid8me .                  # Build image
docker run -it --rm -p 8080:8000 vald8me    # Run & bind to port 8080
```
### Features:
- Facebook, ~~LinkedIn~~ URL validation
- Twitter, Instagram username validation
- Generic/any URL validation
- HTTPS/HTTP detection
	- If a HTTPS/HTTP scheme isn't provided when validating a URL, the server will first try with `https://`. If that fails (perhaps because of a request timeout or an invalid TLS certificate), `http://` will be tried, and the details returned.
	- The `request_url` parameter will be filled with the URL the final request was performed against. Thus, if `https://` succeeded (even if it `404`'ed), you'll be returned an `https://`-prefixed URL.

### Usage:
Available endpoints:
- Generic: `/validate/?url=<ANY URL HERE>`
- Facebook: `/validate/facebook?url=https://www.facebook.com/zuck`
- LinkedIn: `validate/linkedin?url=https://www.linkedin.com/in/satya-nadella-3145136/`
	- LinkedIn validation doesn't really work since LI will nearly always return a `999` status code. However, this could still be useful if you want to know that the URL (probably) points to LinkedIn.
- Twitter: `/validate/twitter?handle=mdo`
- Instagram: `/validate/instagram?username=sundarpichai`

Try from CLI:
```curl -i https://valid8me.stag9.com/validate/twitter\?handle\=parrotmac```

Resp:
```
HTTP/2 200
...
access-control-allow-headers: Content-Type
access-control-allow-methods: GET, HEAD, OPTIONS
access-control-allow-origin: *
...

{"requested_url":"https://twitter.com/parrotmac","status_code":200,"error_message":"","debug_message":""}
```
**Fields**:
- `requested_url`: URL request was ultimately performed against
- `status_code`: Status code from final request
- `error_message`: Relatively user-friendly error message
-`debug_message`: Technical description of a validation issue


### Should I...?
- Be warned? Yes.
- Use it? If you're really careful
- Contribute? If you wanna# internetzValidator / valid8me

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/octet-stream"


file = io.open("IMG_0736.jpeg", "rb")
wrk.body = file:read("*a")

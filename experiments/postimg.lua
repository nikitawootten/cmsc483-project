wrk.method = "POST"
wrk.headers["Content-Type"] = "application/octet-stream"


file = io.open("test_img.jpeg", "rb")
wrk.body = file:read("*a")

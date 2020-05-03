wrk.method = "POST"
wrk.headers["Content-Type"] = "image/jpeg"


file = io.open("test_img.jpeg", "rb")
wrk.body = file:read("*a")

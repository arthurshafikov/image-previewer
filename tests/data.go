package tests

var (
	remoteHostImage1URL         = "http://remote_nginx/test-image-1.jpg"
	remoteHostImage2URL         = "http://remote_nginx/test-image-2.jpeg"
	remoteHostImage3URL         = "http://remote_nginx/test-image-3.jpg"
	remoteHostImageURLExeFile   = "http://remote_nginx/test.exe"
	remoteHostImageURLNotExists = "http://remote_nginx/not-exists-21.jpg"
	remoteHostNotExists         = "http://some-website/not-exists-21.jpg"

	errorHostNotFound           = `{"Data":"couldn't download image from remote host"}`
	errorNotJpgJpegFileResponse = "{\"Data\":\"only jpg and jpeg images are accepted\"}"
)

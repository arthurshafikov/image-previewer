package tests

var (
	remoteHostImage1Url         = "http://remote_nginx/test-image-1.jpg"
	remoteHostImage2Url         = "http://remote_nginx/test-image-2.jpeg"
	remoteHostImage3Url         = "http://remote_nginx/test-image-3.jpg"
	remoteHostImageUrlExeFile   = "http://remote_nginx/test.exe"
	remoteHostImageUrlNotExists = "http://remote_nginx/not-exists-21.jpg"
	remoteHostNotExists         = "http://some-website/not-exists-21.jpg"

	errorHostNotFound = `{"Data":"Get \"` + remoteHostNotExists + `\": dial tcp: ` +
		`lookup some-website: Temporary failure in name resolution"}`
	error404Response            = "{\"Data\":\"remote host has returned: 404 Not Found\"}"
	errorNotJpgJpegFileResponse = "{\"Data\":\"only jpg and jpeg images are accepted\"}"
)

test_evaluate:
		go run ./app/main.go evaluate ./testfile

test_parse:
		go run ./app/main.go parse ./testfile


test_run:
		go run ./app/main.go run ./testfile

test:
	codecrafters test

submit:
	codecrafters submit
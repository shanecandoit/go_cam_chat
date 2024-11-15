
# Linux, macOS
# GOOS=linux  GOARCH=amd64 go build -o webcam_record_linux main.go
# GOOS=darwin GOARCH=amd64 go build -o webcam_record_mac   main.go

# Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o webcam_record.exe main.go

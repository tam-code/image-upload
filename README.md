# Image upload

Image upload is an application built with golang and has rest api to generate upload link and upload image
and async message through kafka to handle some data after upload.

## Installation

Use docker-compose install the app.

```bash
docker-compose up -d
```

## Usage

You can use one of the custom secrets to run the demo:
1. 00000000
2. aaaaaaaa
3. 05f717e5

### Generate upload link
```bash
curl --location 'http://localhost:9521/api/v1/upload-link' \
--header 'X-Secret-Token: 00000000' \
--form 'expiration="2047-10-09T22:50:01.23Z"'
```

### Upload images
```bash
curl --location 'http://localhost:9521/api/v1/images/[UPLOAD-LINK-ID]' \
--header 'Content-Type: multipart/form-data;boundary=AaB03x' \
--form 'images=@"[IMAGE-PATH-FROM-YOUR-MACHINE]"' \
--form 'images=@"[SECOND-IMAGE-PATH-FROM-YOUR-MACHINE]"'
```

### Get image
```bash
curl --location 'http://localhost:9521/api/v1/images/[IMAGE-ID]' \
--header 'Content-Type: application/json'
```

### Get service statistics
```bash
curl --location 'http://localhost:9521/api/v1/statistics' \
--header 'X-Secret-Token: 00000000' \
--header 'Content-Type: application/json'
```

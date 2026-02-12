## API Documentation

# 📚 API Documentation

Base URL: `http://localhost:8080/api/v1`

---

## 1. Convert Image

Mengubah format gambar dari satu format ke format lain (JPG, PNG, WebP).

- **Endpoint**: `/image/convert`
- **Method**: `POST`
- **Content-Type**: `multipart/form-data`

### Request Parameters

| Key            | Type   | Required | Description                                   |
| -------------- | ------ | -------- | --------------------------------------------- |
| `file`         | File   | **Yes**  | File gambar yang akan dikonversi (Max 5MB).   |
| `targetFormat` | String | **Yes**  | Format tujuan. Pilihan: `jpg`, `png`, `webp`. |

### Success Response (200 OK)

```MD
{
    "status": 200,
    "message": "Konversi berhasil",
    "data": {
        "original_file": "my-photo.jpg",
        "result_path": "temp/processed/1709234567-my-photo.png",
        "format": "png"
    }
}
```

## 2. Compression Image

Mengkompresi gambar dengan quality sesuai dengan user

- **Endpoint** : `/image/compress-image`
- **Method** : `POST`
- **Content-Type**: `multipart/form-data`

### Request Parameters

| Key       | Type | Required | Description                          |
| --------- | ---- | -------- | ------------------------------------ |
| `file`    | File | **Yes**  | Tidak ada batas maksimal size gambar |
| `quality` | Int  | **Yes**  | Quality 1 > 100.                     |

### Success Response (200 OK)

```MD
{
    "data": {
        "original_file": "blue-abstract-3840x2160-25119.jpg",
        "quality": 60,
        "result_path": "temp/compressed/1770779587-blue-abstract-3840x2160-25119_compressed.jpeg"
    },
    "message": "Kompresi berhasil",
    "status": 200
}
```

### Error Response (400 Bad Request)

quality tidak berupa int

```MD
{
    "error": "strconv.Atoi: parsing \"60d\": invalid syntax",
    "message": "Quality harus berupa angka",
    "status": 400
}
```

### Error Response (400 Bad Request)

file gambar tidak ditemukan

```MD
{
    "error": "http: no such file",
    "message": "File tidak ditemukan atau tidak valid",
    "status": 400
}
```

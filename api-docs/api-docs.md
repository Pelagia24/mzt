# MZT API Documentation

## Base URL
```
/api/v1
```

## Authentication
All endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your_token>
```

## Error Responses
All endpoints may return the following error responses:

```json
{
    "error": {
        "code": "string",
        "message": "string"
    }
}
```

Common error codes:
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `400 Bad Request` - Invalid request body
- `500 Internal Server Error` - Server error

## Authentication Endpoints

### Sign In
```
POST /auth/signin
```
**Request Body:** `internal/dto.LoginDto`
```json
{
    "email": "string",
    "password": "string"
}
```
**Response:**
```json
{
    "access_token": "string",
    "refresh_token": "string"
}
```

### Sign Up
```
POST /auth/signup
```
**Request Body:** `internal/dto.RegistrationDto`
```json
{
    "name": "string",
    "birthdate": "time.Time",
    "email": "string",
    "phone_number": "string",
    "telegram": "string",
    "city": "string",
    "age": "uint",
    "employment": "string",
    "is_business_owner": "string",
    "position_at_work": "string",
    "month_income": "uint",
    "password": "string"
}
```
**Response:**
```json
{
    "access_token": "string",
    "refresh_token": "string"
}
```

### Refresh Token
```
POST /auth/refresh
```
**Request Body:**
```json
{
    "refresh_token": "string"
}
```
**Response:**
```json
{
    "access_token": "string",
    "refresh_token": "string"
}
```

## User Endpoints

### Get Current User
```
GET /users/me
```
**Response:** `internal/dto.UserInfoDto`
```json
{
    "name": "string",
    "birthdate": "time.Time",
    "email": "string",
    "phone_number": "string",
    "telegram": "string",
    "city": "string",
    "age": "uint",
    "employment": "string",
    "is_business_owner": "string",
    "position_at_work": "string",
    "month_income": "uint"
}
```

### Get Users (Admin Only)
```
GET /users
```
**Response:** Array of `internal/dto.UserInfoAdminDto`
```json
{
    "users": [
        {
            "id": "uuid.UUID",
            "name": "string",
            "birthdate": "time.Time",
            "email": "string",
            "phone_number": "string",
            "telegram": "string",
            "city": "string",
            "age": "uint",
            "employment": "string",
            "is_business_owner": "string",
            "position_at_work": "string",
            "month_income": "uint",
            "course_assignments": "[]CourseDto"
        }
    ]
}
```

### Get User by ID (Admin Only)
```
GET /users/:user_id
```
**Response:** `internal/dto.UserInfoAdminDto`
```json
{
    "id": "uuid.UUID",
    "name": "string",
    "birthdate": "time.Time",
    "email": "string",
    "phone_number": "string",
    "telegram": "string",
    "city": "string",
    "age": "uint",
    "employment": "string",
    "is_business_owner": "string",
    "position_at_work": "string",
    "month_income": "uint",
    "course_assignments": "[]CourseDto"
}
```

### Update User (Admin Only)
```
PUT /users/:user_id
```
**Request Body:** `internal/dto.UpdateUserDto`
```json
{
    "name": "string",
    "birthdate": "time.Time",
    "email": "string",
    "phone_number": "string",
    "telegram": "string",
    "city": "string",
    "age": "uint",
    "employment": "string",
    "is_business_owner": "string",
    "position_at_work": "string",
    "month_income": "uint"
}
```

### Delete User (Admin Only)
```
DELETE /users/:user_id
```

## Course Endpoints

### List Courses
```
GET /courses
```
**Response:** Array of `internal/dto.CourseDto`
```json
{
    "courses": [
        {
            "course_id": "uuid.UUID",
            "name": "string",
            "desc": "string"
        }
    ]
}
```

### Get Course
```
GET /courses/:course_id
```
**Response:** `internal/dto.CourseDto`
```json
{
    "course_id": "uuid.UUID",
    "name": "string",
    "desc": "string"
}
```

### Create Course (Admin Only)
```
POST /courses
```
**Request Body:** `internal/dto.CreateCourseDto`
```json
{
    "name": "string",
    "desc": "string",
    "price": "uint"
}
```

### Update Course (Admin Only)
```
PUT /courses/:course_id
```
**Request Body:** `internal/dto.UpdateCourseDto`
```json
{
    "name": "string",
    "desc": "string",
    "price": "uint"
}
```

### Delete Course (Admin Only)
```
DELETE /courses/:course_id
```

## Lesson Endpoints

### List Lessons
```
GET /courses/:course_id/lessons
```
**Response:** Array of `internal/dto.LessonDto`
```json
{
    "lessons": [
        {
            "lesson_id": "uuid.UUID",
            "course_id": "uuid.UUID",
            "title": "string",
            "desc": "string",
            "video_url": "string",
            "text": "string"
        }
    ]
}
```

### Get Lesson
```
GET /courses/:course_id/lessons/:lesson_id
```
**Response:** `internal/dto.LessonDto`
```json
{
    "lesson_id": "uuid.UUID",
    "course_id": "uuid.UUID",
    "title": "string",
    "desc": "string",
    "video_url": "string",
    "text": "string"
}
```

### Create Lesson (Admin Only)
```
POST /courses/:course_id/lessons
```
**Request Body:** `internal/dto.CreateLessonDto`
```json
{
    "title": "string",
    "desc": "string",
    "video_url": "string",
    "text": "string"
}
```

### Update Lesson (Admin Only)
```
PUT /courses/:course_id/lessons/:lesson_id
```
**Request Body:** `internal/dto.UpdateLessonDto`
```json
{
    "title": "string",
    "desc": "string",
    "video_url": "string",
    "text": "string"
}
```

### Delete Lesson (Admin Only)
```
DELETE /courses/:course_id/lessons/:lesson_id
```

## Course Users Endpoints

### Assign User to Course
```
POST /courses/:course_id/users
```
**Request Body:** `internal/dto.AssignUserToCourseDto`
```json
{
    "user_id": "string"
}
```

### List Users on Course (Admin Only)
```
GET /courses/:course_id/users
```
**Response:** Array of `internal/dto.UserInfoAdminDto`
```json
{
    "users": [
        {
            "id": "uuid.UUID",
            "name": "string",
            "birthdate": "time.Time",
            "email": "string",
            "phone_number": "string",
            "telegram": "string",
            "city": "string",
            "age": "uint",
            "employment": "string",
            "is_business_owner": "string",
            "position_at_work": "string",
            "month_income": "uint",
            "course_assignments": "[]CourseDto"
        }
    ]
}
```

### Remove User from Course (Admin Only)
```
DELETE /courses/:course_id/users/:user_id
```

## Progress Endpoints

### Get Progress
```
GET /courses/:course_id/progress
```
**Response:**
```json
{
    "course_id": "uuid.UUID",
    "progress_percentage": "int"
}
```

### Update Progress
```
PUT /courses/:course_id/progress
```
**Request Body:** `internal/dto.UpdateProgressDto`
```json
{
    "progress": "uint"
}
```

## Notes
1. All timestamps are in ISO 8601 format
2. All UUIDs are in standard UUID v4 format
3. Admin-only endpoints require the user to have admin privileges
4. The progress percentage is calculated based on completed lessons vs total lessons in the course
5. All endpoints are prefixed with `/api/v1`
6. All responses are in JSON format
7. All request bodies should be sent as JSON with `Content-Type: application/json` header
8. All DTOs are defined in the `internal/dto` package

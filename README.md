# IVIS NAS External Web Server

## 개요
- IVIS Lab 외부망에서 동작하는 웹서버.
- 외부망에서 NAS에 접근하기 위한 UI를 제공.

## 개발환경
- Go
- NGINX (Reverse Proxy)
- Redis (세션 유지)
- ZeroTier (NAS와 통신)

## 라우트 정보
- GET /
    - 로그인 페이지
- POST /login
    - 로그인 요청
    - id, pw를 body에 담아서 요청
    - ZeroTier Network (192.168.195.1) 의 /loginCheck 페이지로 요청
    - 성공 시 JWT 토큰 발급, 실패 시 로그인 페이지로 리다이렉트
- GET /logout
    - 로그아웃
    - JWT 토큰 삭제
    - 로그인 페이지로 리다이렉트
    ---
- GET /files/:path
    - 파일 목록 페이지
    - path는 파일 경로
- GET /serve_file/:path
    - 파일 다운로드
    - path는 파일 경로
- GET /download_file/:path (버그 있음)
    - 파일 다운로드
    - path는 파일 경로
- GET /download_dir/:path
    - 디렉토리 다운로드
    - path는 디렉토리 경로
    - zip 파일로 압축하여 다운로드
    ---
- POST /remote_login
    - [IVIS Admin Panel](https://github.com/picel/ivis_admin) 에서 요청하는 아이디/비밀번호 검증 API
    - 성공 시 JWT 토큰 발급, 실패 시 오류 메시지 반환
- POST /remove_token_verify
    - [IVIS Admin Panel](https://github.com/picel/ivis_admin) 에서 요청하는 JWT 토큰 검증 API
    - 성공 시 200 OK, 실패 시 401 Unauthorized 반환


## 실행 화면
- 로그인 페이지
    - ![login](https://user-images.githubusercontent.com/30901178/222891856-9b6833ec-d093-452b-8ebe-2d31ac5d89d3.png)
- 파일 목록 페이지
    - ![files](https://user-images.githubusercontent.com/30901178/222891876-d88f0054-d227-4fab-8341-242232ded8ea.png)
## 주의사항
- key 파일은 git에 올리지 않음.
    - key 파일은 /key 디렉토리에 위치.
    - JWTVerifyKey []byte 정의
example_http
=============

這個名稱之後會更名，主要是拿來練習不使用第三方套件來實作一個 web application，
時做的內容是要作一個大量發送電子報、後台操作、點擊/開信追蹤的服務。

Service
--------

1. webserver: 後台操作，上傳 template、發送名單。template 上傳到 S3、
   發送名單傳送到 SQS。
2. cmd/mailman: 取 SQS 資料回來透過 SES 發送，支援 `docker-compose scale`。 
3. tracker: 點擊/開信追蹤，還在努力中。


To-Do
------

- [All] 支援資料庫紀錄
- [webserver] 建立常用的 sender info. 
- [webserver] template 檔案編檔機制
- [webserver] 登出清除 session
- [tracker] 放入 tracker feeds.

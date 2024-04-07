# 2024 Backend Intern Assignment

 
### 功能
建立和列出廣告
使用 Redis 快取廣告資料以提高效能
使用 MySQL 資料庫儲存廣告資料
 
### 設計
#### 想法和選擇
在設計這個服務時，我們主要考慮了以下幾個方面：
 
**擴充性** 系統應該能夠擴展以支持更多功能和需求。
API 設計
為了提高 API 的易用性，我們採用了以下設計：


使用 Redis 快取廣告資料，以減少資料庫負擔。

使用模組化設計，將 API 和資料庫等功能模組化，方便擴充和維護。
使用資料庫和 Redis 快取來儲存資料，方便擴充資料儲存容量。
## **使用**

### **執行環境設定 (docker-compose.yml)**

這個檔案定義了使用 Docker Compose 執行此廣告投放服務的環境設定。

- **db**：使用 `mysql:latest` 映像作為資料庫服務。
    - **environment**: 設定 MySQL 的環境變數 (僅供開發環境使用，生產環境請勿暴露密碼)
        - MYSQL_ROOT_PASSWORD: MySQL root 使用者的密碼 (請替換為您的密碼)
        - MYSQL_DATABASE: 資料庫名稱 (請替換為您的資料庫名稱)
        - MYSQL_USER: 資料庫使用者的帳號 (請替換為您的帳號)
        - MYSQL_PASSWORD: 資料庫使用者的密碼 (請替換為您的密碼)
    - **ports**: 映射容器埠 3306 到主機埠 3306，使您可以連線到資料庫服務。
    - **volumes**: 可選擇性掛載初始化腳本 (init.sql) 到容器的 `/docker-entrypoint-initdb.d/init.sql` 路徑，在容器啟動時自動執行該腳本 (此為選項)。
        - 也可以選擇性掛載持久化數據的磁碟區 (db-data) 到容器的 `/var/lib/mysql` 路徑，讓資料庫資料可以在容器重新啟動後保存下來 (此為選項)。
    - **restart**: 設定容器重新啟動策略，`unless-stopped` 表示容器除非手動停止，否則會在意外結束後自動重新啟動。
- **redis**：使用 `redis:alpine` 映像作為快取服務。
    - **ports**: 映射容器埠 6379 到主機埠 6379，使您可以連線到快取服務。
    - **restart**: 設定容器重新啟動策略，`unless-stopped` 表示容器除非手動停止，否則會在意外結束後自動重新啟動。
- **go-api** (未啟用)：此範例中沒有啟用 Golang API 容器，您可以加入以下設定啟動 API 服務。
    - **build**: 設定用於建置 Golang API 程式的路徑 (預設為當前目錄)。
    - **ports**: 映射容器埠 8080 到主機埠 8080，使您可以存取 API 服務。
    - **restart**: 設定容器重新啟動策略，`unless-stopped` 表示容器除非手動停止，否則會在意外結束後自動重新啟動。

### **磁碟區**

- **db-data** (未啟用)：此範例中沒有使用持久化數據磁碟區，您可以定義此磁碟區讓資料庫資料可以在容器重新啟動後保存下來 (此為選項)。

### **注意事項**

- 這個範例档中的密碼 (MYSQL_ROOT_PASSWORD, MYSQL_PASSWORD) 僅供開發環境使用，**請勿將其暴露於生產環境**。
- 按照您的實際環境設定資料庫 (MySQL) 的帳號、密碼和資料庫名稱。
- 如果需要使用初始化腳本 (init.sql)，請將其放置在與 docker-compose.yml 檔案相同目錄下，並取消相關行的註解。
- 如果需要使用快取服務 (redis)，請確保 redis 服務正在運行。
- 如果要啟用 Golang API 服務，請取消 `go-api` 區段的註解，並依照您的程式碼設定建置路徑。

### **服務**

### **啟動服務**

您可以使用以下指令啟動使用 Docker Compose 管理的服務:

`docker-compose up -d`

這將會啟動所有定義在 docker-compose.yml 檔案中的服務 (除非有被註解的部份)。

### **停止服務**

您可以使用以下指令停止使用 Docker Compose 管理的服務:

`docker-compose down`

### **連線到資料庫**

您可以使用以下指令連線到 MySQL 資料庫服務:

`mysql -h localhost -P 3306 -u dcard -pi_love_dcard data`

**注意**：請將 `dcard` 替換為您設定的資料庫使用者名稱，並將 `i_love_dcard` 替換為您設定的資料庫使用者密碼。

### 啟動主程式

**1. 安裝 Go**

在遠端電腦上安裝 Go，可從官方網站下載安裝程式：[https://go.dev/dl/](https://go.dev/dl/)。

**2. 傳輸 Go 專案**

**使用 Git:**

- 若您的專案使用 Git 進行版本控制，可使用以下指令將專案倉庫複製到遠端電腦：

`git clone https://your-username@github.com/your-organization/your-project-name.git`

**3. 建置和執行 Go API**

`cd dcard_hw`

**建置 Go API:**

- 使用以下指令建置 Go 二進位檔：

`go build -o main .  # 若您的二進位檔名不同，請用實際名稱替換「main」`

 

 

**執行 Go API:**

- 執行以下指令啟動 API 伺服器：

`./main  # 請用實際的路徑和二進位檔名替換「./main」`

 

## **API**

### **Admin API**

```json
POST /api/v1/ad

Body:

{
  "title": "廣告標題",
  "startAt": "開始時間",
  "endAt": "結束時間",
  "conditions": [
    {
      "ageStart": "年齡下限",
      "ageEnd": "年齡上限",
      "country": ["國家代碼"],
      "platform": ["平台"]
    }
  ]
}
```

### **Public API**

```json
GET /api/v1/ad

Query parameters:

* offset: 偏移量
* limit: 筆數
* age: 年齡
* gender: 性別
* country: 國家代碼
* platform: 平台
```

## **資料庫設計**

為了優化廣告投放服務的效能和彈性，我們採用了以下資料庫設計：

### **資料表**

- **Ads**：此資料表儲存廣告的基本資訊，包括：
    - **id** (INT): 廣告編號 (自動遞增的主鍵)
    - **title** (VARCHAR(255)): 廣告標題 (非空)
    - **start_at** (DATETIME): 廣告投放開始時間 (非空)
    - **end_at** (DATETIME): 廣告投放結束時間 (非空)
    - **min_age** (INT): 廣告投放的最低年齡限制 (預設為 1)
    - **max_age** (INT): 廣告投放的最高年齡限制 (預設為 100)
        
        
- **Conditions**：此資料表儲存廣告的投放條件 (除了年齡以外的條件)，包括：
    - **id** (INT): 條件編號 (自動遞增的主鍵)
    - **type** (ENUM('gender', 'country', 'platform')): 條件的類型 (非空)，包含性別 (gender)、國家 (country) 和平台 (platform)
    - **value** (VARCHAR(255)): 條件的值 (可為空)，例如性別為 'M' 或 'F'、國家為 'TW' 或 'JP' 等
- **AdConditions**：此資料表用於連結廣告和其對應的投放條件，採用多對多的關係。包含：
    - **id** (INT): 關聯編號 (自動遞增的主鍵)
    - **ad_id** (INT): 外鍵，連結至 Ads 資料表的 id (非空)
    - **condition_id** (INT): 外鍵，連結至 Conditions 資料表的 id (非空)

### **索引**

- 在 Ads 資料表上建立一個複合索引 (start_at, end_at)，可以優化按照廣告投放時間查詢的效能。
- 在 Conditions 資料表上建立一個索引 (type)，可以優化按照條件類型查詢的效能。

### **預設值**

- Ads 資料表的 min_age 和 max_age 欄位設定預設值為 1 和 100，方便建立廣告時不必強制指定年齡限制。
- Conditions 資料表預先加入了性別、國家和平台的基本選項，可以未來擴充更多選項。

### **可擴充性**

此資料庫設計允許未來加入更多投放條件類型，只要在 Conditions 資料表新增對應的 type 和 value 即可。同時，AdConditions 資料表采用多對多的關係，可以讓一個廣告包含多個投放條件。

## 程式碼架構

### **目錄**

- **src:** 原始碼
    - **controllers:** 控制器
    - **model:** 資料模型
    - **router:** 路由
    - **services:** 服務
- **test:** 測試
- **docker-compose.yml:** Docker Compose 設定
- **Dockerfile:** Docker 映像檔
- **go.mod:** Go 模組
- **go.sum:** 模組驗證碼
- **init.sql:** 初始化 SQL
- **main.go:** 主程式
- **README.md:** 說明文件

### **檔案說明**

- **src:** 存放專案的原始碼，包含以下子目錄：
    - **controllers:** 處理 API 請求、處理資料和產生回應的控制器檔案。每個控制器檔案通常對應特定 API 資源或端點。
    - **model:** 定義 API 資料結構的資料模型或結構。
    - **router:** 使用 Gin 框架設定 API 路由的程式碼。路由定義了如何將輸入請求映射到控制器中的對應處理程序函數。
    - **services:** 實現業務邏輯或與 API 使用的外部服務（例如資料庫）互動的程式碼。
- **test:** 存放 Go 編寫的單元測試，用於測試應用程式邏輯。
- **main.go:** Go 應用程式的進入點，通常處理初始化任務、路由和啟動 API 伺服器。

### main.go

- **資料庫和 Redis 連線：**
    - `controllers.DBconnect()`：連線資料庫，如果失敗則會中止程式並顯示錯誤（生產環境應更優雅地處理錯誤）。
    - `controllers.Init_redis()`：初始化 Redis 連線，如果失敗則會中止程式（建議同上）。
    - `go func() { ... }()`：在一個新的 Go 程式（goroutine）中，再次呼叫 `controllers.DBconnect()` 和 `controllers.Init_redis()`，確保連線成功。
- **創建路由器：**
    - `router := gin.Default()`：使用 Gin 框架創建一個預設路由器，用來處理 API 請求和映射到對應的處理函數。
- **定義 API 路由：**
    - `router.POST("api/v1/ad", service.CreateAd)`：新增一個 POST 請求路由，URL 為 `/api/v1/ad`，由 `service.CreateAd` 函數處理，用於建立廣告。
    - `router.GET("api/v1/ad", service.ListAds)`：新增一個 GET 請求路由，URL 為 `/api/v1/ad`，由 `service.ListAds` 函數處理，用於列出廣告。
- **啟動伺服器：**
    - `router.Run(":8080")`：啟動 Web 伺服器，監聽 8080 埠，開始接收和處理 API 請求。

**概括來說，這段程式碼建立了一個簡單的 Web API，使用 Gin 框架來處理 API 請求，並提供以下功能：**

- 連線資料庫和 Redis。
- 提供創建廣告的 API 端點（POST /api/v1/ad）。
- 提供列出廣告的 API 端點（GET /api/v1/ad）。

## model

**資料結構：**

- **Ad（廣告）結構體：**
    - Title（標題）：string
    - StartAt（開始時間）：time.Time
    - EndAt（結束時間）：time.Time
    - Conditions（條件）：[]Condition
- **Condition（條件）結構體：**
    - AgeStart（起始年齡）：int
    - AgeEnd（結束年齡）：int
    - Gender（性別）：[]string
    - Country（國家）：[]string
    - Platform（平台）：[]string
- **Search_Condition（搜索條件）結構體：**
    - Age（年齡）：int
    - Gender（性別）：string
    - Country（國家）：string
    - Platform（平台）：string
    - Limit（限制數量）：int
    - Offset（偏移量）：int
- **Result（結果）結構體：**
    - Title（標題）：string
    - EndAt（結束時間）：time.Time

**全域變數：**

- All_platform：所有平台的列表（["android", "ios", "web"]）
- All_gender：所有性別的列表（["M", "F"]）

## service.ad

**CreateAd 函式：**

- 處理建立廣告請求（POST /api/v1/ad）。
- 流程：
    1. 從請求體解析 JSON 資料，繫結到 `Ad` 結構體。
    2. 驗證 `Ad` 資料（標題、開始時間、結束時間等）。
    3. 補充 `Ad` 資料中的缺省值（國家、性別、平台）。
    4. 呼叫 `controllers.Create_ad` 函式，將 `Ad` 資料存入資料庫。
    5. 返回 201 Created 狀態碼。

**ListAds 函式：**

- 處理列出廣告請求（GET /api/v1/ad）。
- 流程：
    1. 從查詢參數（Query Params）：
        - 提取過濾條件（年齡、性別、國家、平台）。
        - 提取分頁資訊（Offset、Limit，並做驗證）。
    2. 建立 `Search_Condition` 結構體，存放所有查詢條件。
    3. 呼叫 `controllers.Find_ad` 函式，根據條件和分頁從資料庫獲取廣告。
    4. 如果獲取成功，返回 200 OK 狀態碼，並以 JSON 格式返回廣告列表。
    5. 如果失敗，返回 500 內部伺服器錯誤。

## controller.contoller_ad

### Find

**函式的目的：**

- `Find_ad` 函式用於搜索符合指定條件的廣告。
- 它優先從 Redis 緩存中獲取資料，以提高響應速度。
- 如果 Redis 失敗或資料不存在，會轉而查詢資料庫。

**函式參數：**

- `condition (model.Search_Condition)`：搜索條件，包含過濾條件和分頁資訊。

**函式返回值：**

- `[]model.Result (nil on error)`：匹配的廣告結果陣列，如果出錯則返回空值。
- `error`：搜索過程中遇到的任何錯誤。

**函式流程：**

1. **嘗試從 Redis 獲取廣告：**
    - 呼叫 `find_adByredis` 函式嘗試從 Redis 獲取廣告。
    - 如果成功（`ads != nil`），則表示在 Redis 中找到了結果，直接返回。
    - 如果失敗（`err != nil`），則報錯並返回。
2. **如果 Redis 失敗，則查詢資料庫：**
    - 記錄 "redis fail" 訊息。
    - 呼叫 `find_adBysql` 函式查詢資料庫。
    - 如果查詢成功，則返回資料庫中的結果。
    - 如果查詢失敗，則報錯並返回。
3. **應用分頁限制：**
    - 呼叫 `getSlice` 函式，根據 `condition.Limit` 和 `condition.Offset` 對結果進行分頁。
    - 返回分頁後的結果和錯誤（如果有）。
4. sql 
    
    ```sql
    SELECT a.id, a.title, a.start_at, a.end_at
    FROM Ads a 
    INNER JOIN AdConditions ac ON a.id = ac.ad_id 
    INNER JOIN Conditions c ON ac.condition_id = c.id
    
    WHERE
     a.start_at <= NOW() AND a.end_at >= NOW() 
    AND 
    (a.min_age <= 24 AND a.max_age >= 24)
    AND 
    
    ((c.type = 'gender' AND c.value = ('F')) 
    OR (c.type = 'country' AND c.value = ('TW')) 
    OR (c.type = 'platform' AND c.value = ('ios')) )
    
    group by ad_id 
    having  COUNT(ad_id)=3
    ORDER by end_at ASC
    ```
    
    這是廣告搜尋的 SQL 查詢語句，用於在資料庫中找出符合指定條件的廣告。
    
    **資料表：**
    
    - Ads：儲存廣告基本資料，例如廣告 ID (id)、標題 (title)、開始時間 (start_at)、結束時間 (end_at)、最小年齡 (min_age)、最大年齡 (max_age) 等。
    - AdConditions：儲存廣告和條件的關聯，例如廣告條件 ID (ad_id) 和條件 ID (condition_id)。
    - Conditions：儲存廣告條件的詳細資料，例如條件類型 (type)、條件值 (value)。
    
    **查詢條件：**
    
    1. **過濾當前有效的廣告：**
        - `a.start_at <= NOW()`: 篩選開始時間小於或等於當前時間的廣告。
        - `a.end_at >= NOW()`: 篩選結束時間大於或等於當前時間的廣告。
    2. **年齡限制 (可選)：**
        - `(a.min_age <= 24 AND a.max_age >= 24)`: 篩選最小年齡小於或等於 24 且最大年齡大於或等於 24 的廣告，表示此廣告僅針對 24 歲的使用者。
    3. **多重條件篩選：**
        - `((c.type = 'gender' AND c.value = ('F'))`: 條件一：篩選條件類型為 'gender' (性別) 且值為 'F' (女性) 的廣告。
        - `OR (c.type = 'country' AND c.value = ('TW'))`: 條件二：篩選條件類型為 'country' (國家) 且值為 'TW' (台灣) 的廣告。
        - `OR (c.type = 'platform' AND c.value = ('ios'))`: 條件三：篩選條件類型為 'platform' (平台) 且值為 'ios' 的廣告。
            - 使用 OR 運算子表示符合下列任何一個條件的廣告都會被篩選出來。
    4. **分組和排序：**
        - `group by ad_id`: 根據廣告 ID (ad_id) 分組，將具有相同廣告 ID 的資料歸納在一起。
        - `having COUNT(*) = 3`: 篩選分組後資料行數等於 3 的群組，也就是說僅保留符合所有三個條件 (當前時段、年齡限制、廣告屬性) 的廣告。
        - `ORDER BY end_at ASC`: 根據結束時間 (end_at) 遞增排序，即優先列出即將結束的廣告。
    
    **總而言之，這段查詢會過濾出當前時段內有效的廣告，並根據可選的年齡限制以及性別、國家、平台等屬性進行篩選，最終按照結束時間排序列出符合所有條件的廣告。**
    

### **Create**

 

**Create_ad 函式的目的：**

- 在資料庫中建立新的廣告記錄和相關的條件。
- 接受一個 `Ad` 結構體作為參數，其中包含廣告的詳細資訊。
- 如果建立過程中有錯誤，則返回錯誤。

**Create_ad 函式的流程：**

1. **開啟交易：** 使用 `Db.Begin()` 開啟一個資料庫交易，確保資料一致性。
2. **插入廣告：**
    - 根據是否有年齡限制，使用不同的 SQL 語句插入廣告記錄。
    - 如果插入失敗，則回滾交易並返回錯誤。
3. **獲取廣告 ID：** 使用 `Db.Raw` 查詢新建廣告的 ID。
4. **建立廣告條件：** 呼叫 `Create_condition` 函式，建立廣告和條件之間的關聯。
5. **提交或回滾：** 如果所有操作都成功，則提交交易並刷新廣告緩存。否則回滾交易。

**Create_condition 函式的目的：**

- 在資料庫中建立廣告和條件之間的關聯。
- 接受一個 `Condition` 結構體和廣告 ID 作為參數。
- 如果建立過程中有錯誤，則返回錯誤。

**Create_condition 函式的流程：**

1. **並發獲取子句：**
    - 使用 `sync.WaitGroup` 並發執行 `get_clause` 函式，加速獲取國家、平台、性別條件的子句。
    - 等待所有並發任務完成。
2. **構建 SQL 查詢：** 動態構建 SQL 查詢語句，把廣告 ID 和條件 ID 插入到 `AdConditions` 表中。
3. **執行查詢：** 執行 SQL 查詢，建立廣告條件記錄。
4. **返回結果：** 如果所有操作都成功，則返回 `nil`。否則返回錯誤。
5. sql
    
    ```sql
    -- Replace `<ad_id>` with the actual ad ID (retrieved in step 2 or your existing ad)
    INSERT INTO AdConditions (ad_id, condition_id)
    SELECT <ad_id>, c.id
    FROM Conditions c
    WHERE  
      OR c.type = 'gender' AND c.value IN ('M', 'F')  -- Target both genders
      OR c.type = 'country' AND c.value IN ('TW', 'JP')  -- Target both countries (Taiwan & Japan)
      OR c.type = 'platform' AND c.value IN ('android', 'ios', 'web');  -- Target all platforms
    
    ```
    
    - `INSERT INTO AdConditions (ad_id, condition_id)`: 指示要將資料插入到 `AdConditions` 資料表，並指定要插入的欄位（`ad_id` 和 `condition_id`）。
    - `SELECT 1, c.id`: 這個子查詢會選擇一個常數值 1 和資料表 `Conditions` 中的 `id` 欄位。
    - `FROM Conditions c`: 指定要從資料表 `Conditions` 取得資料。
    - `WHERE`: 後面的條件用於篩選要插入的資料。
        - `OR`: 表示下列條件是 OR 邏輯，只要符合其中之一就會被加入。
            - `c.type = 'gender' AND c.value IN ('M', 'F')`: 條件一：篩選條件類型為 'gender' (性別) 且值為 'M' (男性) 或 'F' (女性) 的資料。
            - `c.type = 'country' AND c.value IN ('TW', 'JP')`: 條件二：篩選條件類型為 'country' (國家) 且值為 'TW' (台灣) 或 'JP' (日本) 的資料。
            - `c.type = 'platform' AND c.value IN ('android', 'ios', 'web')`: 條件三：篩選條件類型為 'platform' (平台) 且值為 'android'、'ios' 或 'web' 的資料。
    
    **範例中的註解：**
    
    - `<ad_id>`: 這是一個佔位符，需要替换成實際的廣告 ID (應該在程式碼的其他部分取得)。
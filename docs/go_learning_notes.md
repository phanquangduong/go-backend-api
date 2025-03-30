- Folder internal chứ mã cục bộ của dự án

- pkg package sử dụng chung trong các dự án nằm trong source Go

- gin là 1 framework viết bằng Go

- go.mod:

  - indirect những thành phần không phụ thuộc vào chính dự án

- 1 trong thế mạnh của go là struct

- Dấu \* là con trỏ và dấu & là địa chỉ con trỏ:

  - type UserController struct {
    }

  - func NewUserController() \*UserController {
    return &UserController{}
    }

- Mục đích khi sử dụng con trỏ để quản lý bộ nhớ hiệu quả hơn và đảm bảo các thay đổi trên struct được phản ánh ở mọi nơi sử dụng nó. Nếu không sử dụng con trỏ chương trình sẽ chạy mà không gặp lỗi tuy nhiên mỗi lời gọi controller.NewUserController(), Go sẽ tạo 1 bản sao mới của struct UserController gây hệ quả:

  - Không hiệu quả về bộ nhớ
  - Không thể dữ trạng thái: Nếu UserController có các trường lưu trữ dữ liệu, các thay đổi trên bản sao sẽ không được phản ánh về bản gốc

- Định nghĩa kiểu struct UserController:

  -type UserController struct {
  userService \*service.UserService
  }

  - Ở đây, bạn đang định nghĩa một struct có tên là UserController.
    Trong struct này, có một trường (field) có tên là userService, và kiểu của trường này là con trỏ đến \*service.UserService. Điều này có nghĩa là UserController sẽ có một sự phụ thuộc (dependency) vào một instance của UserService.

- Khởi tạo và gán giá trị cho userService:

  - func NewUserController() \*UserController {
    return &UserController{
    userService: service.NewUserService(),
    }
    }
  - Hàm NewUserController sẽ tạo và trả về một instance mới của UserController.
    Trong khi khởi tạo UserController, trường userService sẽ được gán giá trị là kết quả trả về từ service.NewUserService() — tức là một instance của UserService sẽ được tạo và gán vào userService.

- Gin vs Logs Handler

  - Có 2 loại là sugar và logger
  - Khi nào dùng sugar và khi nào dùng logger:

- Viper là một thư viện Go mạnh mẽ để quản lý cấu hình. Nó hỗ trợ nhiều loại định dạng file cấu hình như JSON, YAML, TOML, HCL, INI, và cũng có thể lấy cấu hình từ biến môi trường, cờ dòng lệnh, hay thậm chí là một cơ sở dữ liệu. Viper giúp làm cho việc quản lý cấu hình dễ dàng hơn trong các ứng dụng Go.

- Dùng require để test nếu faild các hàm bên dưới không thực thi
- Dùng assert để test nếu faild các hàm bên dưới vẫn thực thi

- DUMP Database:
  docker exec mydb mysqldump -uroot -p12345 --databases shopGO --add-drop-database --add-drop-table --add-drop-trigger --add-locks --no-data > migrations/shopGO.sql

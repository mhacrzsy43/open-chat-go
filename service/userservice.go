package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// @GetUserList
// @Summary 查询用户
// @Tags 用户模块
// @Success 200 {string} json {"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //正常 -1 失败
		"message": data,
	})
}

// @Summary 登陆
// @Tags 用户模块
// @Accept  json
// @Produce  json
// @Param   user  body  CreateUserRequest  true  "用户注册信息"
// @Success 200 {string} json {"code", "message"}
// @Router /user/login [POST]
func GetUserByName(c *gin.Context) {

	var req struct {
		Name       string `json:"username"`
		Password   string `json:"password"`
		Repassword string `json:"repassword"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式不正确"})
		return
	}

	if req.Name == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    0, //正常 -1 失败
			"message": "用户名或密码不能为空",
		})
		return
	}

	user := models.FindUserByName(req.Name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //正常 -1 失败
			"message": "用户不存在",
		})
		return
	}
	flag := utils.ValidPassword(req.Password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //正常 -1 失败
			"message": "密码不正确",
		})
		return
	}
	data := models.FindUserByNameAndPwd(req.Name, user.PassWord)
	utils.RespOK(c.Writer, data)
}

type UpdateUserRequest struct {
	ID       int32  `json:"id" binding:"required"`        // 用户ID, 必填
	Name     string `json:"name" binding:"omitempty"`     // 用户名, 可选
	Email    string `json:"email" binding:"omitempty"`    // 邮箱, 可选
	Phone    string `json:"phone" binding:"omitempty"`    // 电话, 可选
	PassWord string `json:"password" binding:"omitempty"` // 电话, 可选
}

// CreateUserRequest represents the request payload for user creation.
type CreateUserRequest struct {
	Name       string `json:"name" binding:"required"`       // 用户名
	Password   string `json:"password" binding:"required"`   // 密码
	Repassword string `json:"repassword" binding:"required"` // 确认密码
}

// RegisterUser creates a new user
// @Summary 注册用户
// @Description 注册新用户
// @Tags 用户模块
// @Accept  json
// @Produce  json
// @Param   user  body  CreateUserRequest  true  "用户注册信息"
// @Success 200 {string} json {"code", "message"}
// @Router /user/register [post]
func RegisterUser(c *gin.Context) {
	var req struct {
		Name       string `json:"name"`
		Password   string `json:"password"`
		Repassword string `json:"repassword"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式不正确"})
		return
	}

	if req.Name == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    0, //正常 -1 失败
			"message": "用户名或密码不能为空",
		})
		return
	}

	if data := models.FindUserByName(req.Name); data.Name != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "用户名已注册",
		})
		return
	}

	if req.Password != req.Repassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1, //正常 -1 失败
			"message": "两次密码不一致",
		})
		return
	}
	salt := rand.Int31()
	user := models.UserBasic{
		Name:     req.Name,
		PassWord: utils.MakePassword(req.Password, fmt.Sprintf("%06d", salt)),
		Salt:     salt,
	}

	models.CreateUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0, //正常 -1 失败
		"message": "注册用户成功",
	})
}

// @DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json {"code", "message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"message": "删除用户成功",
	})
}

// UpdateUser 修改用户
// @Summary 修改用户
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param user body UpdateUserRequest true "用户信息"
// @Success 200 {string} json {"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	var user models.UserBasic
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数错误",
		})
		return
	}

	fmt.Println("用户信息：", user)

	// 数据校验
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "修改参数格式不对",
		})
		return
	}

	// 调用模型层的更新用户函数
	models.UpdateUser(user)

	c.JSON(http.StatusOK, gin.H{
		"message": "修改用户成功",
	})
}

// UpdatePassword 修改用户
// @Summary 修改密码
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param user body UpdateUserRequest true "用户信息"
// @Success 200 {string} json {"code", "message"}
// @Router /user/updatePassword [post]
func UpdatePassword(c *gin.Context) {
	var user models.UserBasic
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数错误",
		})
		return
	}
	data := models.FindUserByName(user.Name)
	if data.Name == "" {
		c.JSON(200, gin.H{
			"message": "用户不存在",
		})
		return
	}

	fmt.Println("用户信息：", user)
	data.PassWord = utils.MakePassword(user.PassWord, fmt.Sprintf("%06d", data.Salt))
	models.UpdateUser(data)
	c.JSON(http.StatusOK, gin.H{
		"message": "修改用户成功",
	})
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	// 定义消息处理函数
	handleMessage := func(msg string) {
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]: %s", tm, msg)
		err := ws.WriteMessage(websocket.TextMessage, []byte(m))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}

	// 在新的goroutine中订阅消息
	go utils.Subscribe(c, utils.PublishKey, handleMessage)

	// 保持函数活跃以维持WebSocket连接
	// 这里可以根据需要添加逻辑来处理WebSocket连接的关闭等
	select {
	case <-c.Done(): // 当Gin上下文结束时停止循环
		fmt.Println("MsgHandler: Gin context is done")
		return
	}
}
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// SearchFriends 查找好友
// @Summary 查找好友
// @Description 根据用户ID查找其好友列表
// @Tags 用户模块
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param userId formData string true "用户ID"
// @Success 200 {object} map[string]interface{} "返回好友列表和成功消息"
// @Failure 400 {object} map[string]interface{} "返回错误消息"
// @Router /user/getFriends [post]
func SearchFriends(c *gin.Context) {
	token, ok := c.Value("token").(string)
	if !ok {
		// 如果Token不存在，返回错误
		utils.RespFail(c.Writer, "token无效")
		return
	}
	friends := models.SearchFriend(token)
	utils.RespOKList(c.Writer, friends, len(friends))
}

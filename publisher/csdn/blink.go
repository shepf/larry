package csdn

import (
	"bytes"
	"fmt"
	"github.com/dghubble/sling"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// 要提交的参数  Form Body 格式 form-data 格式
// BlinkStatusUpdateParams represents a CSNS Blink, previously called a Blink status.
type BlinkStatusUpdateParams struct {
	Content    string `url:"content,omitempty"`
	ActivityId string `url:"activityId,omitempty"`
}

// BlinkService provides a method for account credential verification.
type BlinkService struct {
	sling *sling.Sling
}

// newBlinkService returns a new BlinkService.
func newBlinkService(sling *sling.Sling) *BlinkService {
	return &BlinkService{
		sling: sling.Path("blink/"),
	}
}

// 更新csdn blink 状态
// Update updates the user's blink status, also known as Blink.
func (s *BlinkService) Update(status string, params *BlinkStatusUpdateParams) (*BlinkStatusUpdateParams, *http.Response, error) {
	if params == nil {
		params = &BlinkStatusUpdateParams{}
	}
	params.Content = status
	csdn := new(BlinkStatusUpdateParams)
	apiError := new(APIError)

	//log.Printf("csdn blink 发布内容： %v", status)
	//activityId 是 主题活动id，例如 要在主题活动 每日学习打卡 是 93；开源项目推荐 是 47
	extraParams := map[string]string{
		"content":    status,
		"activityId": "47",
	}
	request := newMultipartRequest("https://blink-open-api.csdn.net/v1/pc/blink/sendBlink", extraParams)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("client.Do err:")
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		//fmt.Println(resp.Header)
		//fmt.Println(body)
	}

	//log.Printf("csdn blink 请求路径： %v", req.URL)
	//log.Printf("csdn blink 请求方法： %v", req.Method)
	//log.Printf("csdn blink 请求Content-Type： %v", req.Header.Get("Content-Type"))
	//
	//
	//resp, err := s.sling.Do(req, nil,apiError)

	if 200 != resp.StatusCode {
		log.Fatalf("csdn blink 发布失败 %v", resp)
	}

	return csdn, resp, relevantError(err, *apiError)
}

// for multipart/form-data;
type formDataBodyProvider struct {
	payload interface{}
}

func (p formDataBodyProvider) ContentType() string {
	return "multipart/form-data"
}

func (p formDataBodyProvider) Body() (io.Reader, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	extraParams := map[string]string{
		"content":    "xxxxxxx",
		"activityId": "",
	}
	for key, val := range extraParams {
		_ = writer.WriteField(key, val)
	}
	writer.Close()

	reader := bytes.NewReader(body.Bytes())

	return reader, nil

}

func newMultipartRequest(url string, params map[string]string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	str, _ := os.Getwd()
	fmt.Println("获取文件cookie.txt路径: ", str)
	//因为cookie可能会变化，这里从文件读取，每次发请求时，相当取一次。cookie更新时，更新文件即可~!
	//使用ioutil.ReadFile 直接从文件读取到 []byte中
	f, err := ioutil.ReadFile("cookie.txt")
	if err != nil {
		fmt.Errorf("从文件 cookie.txt 获取cookie 失败！")
		panic(err)
	}
	fmt.Println(string(f))
	req.Header.Set("Cookie", string(f))

	return req
}

package main

import (
	"encoding/json"
	_ "encoding/json"
	_ "errors"
	"fmt"
	_ "fmt"
	"github.com/aiomonitors/godiscord"
	_ "github.com/aiomonitors/godiscord"
	"io/ioutil"
	"net/http"
	_ "net/http"
	"time"
)

type Product struct {
	MainItem MainItem `json:"MainItem"`
	Sku string
}

type MainItem struct {
	Description Description `json:"Description"`
	Image Image `json:"Image"`
	InStock bool `json:"InStock"`
	Price float64 `json:"FinalPrice"`
}

type Image struct {
	Normal Normal `json:"Normal"`
}

type Description struct {
	Title string `json:"Title"`
	UrlKeywords string `json:"UrlKeywords"`
}

type Normal struct {
	ImageName string `json:"ImageName"`
}

func sendWebhook(product Product, discordWebhook string)  {
	imageUrl := fmt.Sprintf("https://c1.neweggimages.com/ProductImageCompressAll1280/%s", product.MainItem.Image.Normal.ImageName)
	productUrl := fmt.Sprintf("https://www.newegg.com/%s/p/%s", product.MainItem.Description.UrlKeywords, product.Sku)
	emb := godiscord.NewEmbed(product.MainItem.Description.Title, "", productUrl)

	emb.SetThumbnail(imageUrl)
	emb.SetColor("#e7a348")
	emb.SetFooter("NewEgg Monitor by Brennan", "")
	emb.AddField("Price",  fmt.Sprintf("$%.2f", product.MainItem.Price), true)
	emb.AddField("SKU", product.Sku, true)

	emb.Username = "NewEgg Monitor"

	emb.SendToWebhook(discordWebhook)
}

func newMonitor()  {
	client := &http.Client{}

	discordWebhook := "web-hook"
	sku := "N82E16819113567"

	request, requestError := http.NewRequest("GET", 	fmt.Sprintf("https://www.newegg.com/product/api/ProductRealtime?ItemNumber=%s",sku), nil)
	if requestError != nil {
		fmt.Println(requestError)
		return
	}
	request.Close = true

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"cache-control": "no-store,no-cache,must-revalidate,proxy-revalidate,max-age=0",
		"pragma": "no-cache",
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, responseError := client.Do(request)
	if responseError != nil {
		fmt.Println(responseError)
		return
	}

	var product Product

	body, bodyErr := ioutil.ReadAll(response.Body)

	if bodyErr != nil {
		fmt.Println(bodyErr)
		return
	}
	response.Body.Close()

	err := json.Unmarshal(body, &product)

	if err != nil {
		fmt.Println(err)
		return
	}

	product.Sku = sku

	stock := product.MainItem.InStock
	if stock {
		sendWebhook(product, discordWebhook)
	}
}

func main() {
	i := true
	for i == true {
		time.Sleep(750 * time.Millisecond)

		newMonitor()
	}
}


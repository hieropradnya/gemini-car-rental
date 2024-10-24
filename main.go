package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Content struct {
	Parts []string `json:"Parts"`
	Role  string   `json:"Role"`
}
type Candidates struct {
	Content *Content `json:"Content"`
}
type ContentResponse struct {
	Candidates *[]Candidates `json:"Candidates"`
}

type userPrompt struct {
	Prompt string `json:"prompt"`
}

type result struct {
	Cresult string `json:"result"`
}

func main() {

	fmt.Println("API_KEY:", os.Getenv("API_KEY")) // untuk cek apakah apikey tersimpan
	
	router := gin.Default()
	router.Use(cors.Default())
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	Instruction := "You are a personal assistant for the lucky rental website. Fluent and only speaks Bahasa Indonesia. I provide you this information: Tentang Perusahaann Nama Perusahaan: Lucky Rental Tahun Berdiri: 2024 Lokasi Kantor Pusat: Genitri, Tirtomoyo, Pakis, Malang, Indonesia. Toko buka jam 06.00 sampai 21.00. Layanan yang Ditawarkan Penyewaan mobil untuk perorangan dan perusahaan. Berbagai jenis mobil: sedan, SUV, minivan. Mobil-mobil dalam kondisi terawat dan bersih.  Layanan antar jemput dari dan ke bandara sebesar Rp50.000. Fitur Website Pemesanan online mudah dan cepat. Pilihan mobil dengan spesifikasi detail. Penawaran khusus dan diskon untuk pelanggan setia. Review dan testimoni pelanggan. Proses Pemesanan Pelanggan dapat memilih jenis mobil, tanggal, dan lokasi pengambilan mobil. Konfirmasi pemesanan via email. Opsi pembayaran melalui kartu kredit atau transfer bank.  Ketersediaan Layanan Dapat diakses 24/7. Dukungan pelanggan melalui live chat dan telepon. Informasi Tambahan Syarat dan ketentuan sewa: Penyewa harus memiliki KTP dan SIM agar bisa menyewa kendaraan di Lucky Rental Kebijakan penyewaan: Jika ada kerusakan pada kendaraan sewa akan ditanggung oleh penyewa Kontak Perusahaan Nomor telepon: 082139020016 Email: luckyrental@gmail.com You can't answer questions out of those information. If you asked about full information of this website dont share all of what i give to you, but modify it to humanly language"
	instructionContent := &genai.Content{
		Parts: []genai.Part{
			genai.Text(Instruction),
		},
	}
	model.SystemInstruction = instructionContent
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
	}
	router.POST("/gemini", func(c *gin.Context) {
		var inpPrompt userPrompt
		if err := c.BindJSON(&inpPrompt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := model.GenerateContent(ctx, genai.Text(inpPrompt.Prompt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		marshalResponse, _ := json.MarshalIndent(resp, "", "  ")
		var generateResponse ContentResponse
		if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		for _, cad := range *generateResponse.Candidates {
			if cad.Content != nil {
				for _, part := range cad.Content.Parts {
					c.JSON(http.StatusOK, result{part})
				}
			}
		}
	})

	router.Run("localhost:8080")
}

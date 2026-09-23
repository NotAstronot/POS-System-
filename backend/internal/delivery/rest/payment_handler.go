package rest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type QRISHandler struct{}

func NewQRISHandler() *QRISHandler {
	return &QRISHandler{}
}

func crc16CCITT(data string) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range []byte(data) {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 != 0 {
				crc = (crc >> 1) ^ 0x8408
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

func tlv(tag string, value string) string {
	return tag + fmt.Sprintf("%02d", len(value)) + value
}

func buildQRISPayload(amount float64, merchantName, city string) string {
	guid := "ID.CO.QRIS.WWW"
	mid := "09876543210"
	merchant := tlv("00", guid) + tlv("01", mid)

	amountStr := fmt.Sprintf("%.2f", amount)
	amountStr = strings.Replace(amountStr, ".", "", 1)

	payload := "000201010211" +
		tlv("26", merchant) +
		tlv("52", "5361") +
		tlv("53", "360") +
		tlv("58", "ID") +
		tlv("54", amountStr) +
		tlv("59", merchantName) +
		tlv("60", city) +
		"6304"

	crc := crc16CCITT(payload)
	return payload + fmt.Sprintf("%04X", crc)
}

func (h *QRISHandler) Generate(c *gin.Context) {
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount harus lebih dari 0"})
		return
	}

	payload := buildQRISPayload(req.Amount, "POS MERCHANT", "JAKARTA")

	png, err := qrcode.Encode(payload, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat QR code"})
		return
	}

	var buf bytes.Buffer
	buf.Write(png)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"qris":      payload,
			"qr_base64": base64.StdEncoding.EncodeToString(buf.Bytes()),
		},
	})
}

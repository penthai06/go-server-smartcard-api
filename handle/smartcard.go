package handle

import (
	"os"
	"strings"

	"github.com/Napat/sscard"
	"github.com/ebfe/scard"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/color"
	"github.com/varokas/tis620"
)

type ReponseData struct {
	CID        string `json:"cid"`
	FullNameTh string `json:"full_name_th"`
	FullNameEn string `json:"full_name_en"`
	Birth      string `json:"birth"`
	Gender     string `json:"gender"`
	IssueExp   string `json:"issue_exp"`
	Image      string `json:"image"`
}

func HandleReader(c *fiber.Ctx) (ReponseData, error) {
	// Establish a PC/SC context
	context, err := scard.EstablishContext()
	if err != nil {
		color.Error.Println("Error EstablishContext:", err)
		return ReponseData{}, err
	}

	// Release the PC/SC context (when needed)
	defer context.Release()

	// List available readers
	readers, err := context.ListReaders()
	if err != nil {
		color.Error.Println("Error ListReaders:", err)
		return ReponseData{}, err
	}

	// Use the first reader
	reader := readers[0]
	color.Warn.Println("Using reader:", reader)

	// Connect to the card
	card, err := context.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		color.Error.Println("Error Connect:", err)
		return ReponseData{}, err
	}

	// Disconnect (when needed)
	defer card.Disconnect(scard.LeaveCard)

	// Send select APDU
	selectRsp, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardSelect)
	if err != nil {
		color.Error.Println("Error Transmit:", err)
		return ReponseData{}, err
	}
	color.Warn.Println("resp sscard.APDUThaiIDCardSelect: ", selectRsp)

	cid, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardCID)
	if err != nil {
		color.Error.Println("Error APDUGetRsp: ", err)
		return ReponseData{}, err
	}
	color.Info.Printf("cid: %s\n", string(cid))
	cidTrim := strings.Trim(string(cid), "\u0000\u0000")

	fullnameEN, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardFullnameEn)
	if err != nil {
		color.Error.Println("Error APDUGetRsp: ", err)
		return ReponseData{}, err
	}
	// color.Info.Printf("fullnameEN: %s\n", string(fullnameEN))

	fullnameTH, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardFullnameTh)
	if err != nil {
		color.Error.Println("Error APDUGetRsp: ", err)
		return ReponseData{}, err
	}
	// color.Info.Printf("fullnameTH: %s\n", string(fullnameTH))
	fullnameThUtf8 := string(tis620.ToUTF8(fullnameTH))
	fullnameTrim1 := strings.Trim(fullnameThUtf8, "\u0000\u0000")
	fullnameTrim1 = strings.Trim(fullnameTrim1, " ")

	birth, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardBirth)
	if err != nil {
		color.Error.Println("Error APDUGetRsp: ", err)
		return ReponseData{}, err
	}
	// color.Info.Printf("birth: %s\n", string(birth))
	birthTrim := strings.Trim(string(birth), "\u0000\u0000")

	gender, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardGender)
	if err != nil {
		color.Error.Println("Error APDUGetRsp: ", err)
		return ReponseData{}, err
	}
	// color.Info.Printf("gender: %s\n", string(gender))
	genderTrim := strings.Trim(string(gender), "\u0000\u0000")

	// issuer, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardIssuer)
	// if err != nil {
	// 	color.Error.Println("Error APDUGetRsp: ", err)
	// 	return ReponseData{}, err
	// }
	// color.Info.Printf("issuer: %s\n", string(issuer))

	// issueDate, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardIssuedate)
	// if err != nil {
	// 	color.Error.Println("Error APDUGetRsp: ", err)
	// 	return ReponseData{}, err
	// }
	// color.Info.Printf("issueDate: %s\n", string(issueDate))

	issueExp, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardExpiredate)
	if err != nil {
		color.Error.Println("Error APDUGetRsp: ", err)
		return ReponseData{}, err
	}
	// color.Info.Printf("issueExp: %s\n", string(issueExp))
	issueExpTrim := strings.Trim(string(issueExp), "\u0000\u0000")

	// address, err := sscard.APDUGetRsp(card, sscard.APDUThaiIDCardAddress)
	// if err != nil {
	// 	color.Error.Println("Error APDUGetRsp: ", err)
	// 	return ReponseData{}, err
	// }
	// color.Info.Printf("address: %s\n", string(address))

	cardPhotoJpg, err := sscard.APDUGetBlockRsp(card, sscard.APDUThaiIDCardPhoto, sscard.APDUThaiIDCardPhotoRsp)
	if err != nil {
		color.Error.Println("Error: ", err)
		return ReponseData{}, err
	}
	// fmt.Printf("Image binary: % 2X\n", cardPhotoJpg)

	err = os.MkdirAll("./imgs", os.ModePerm)
	if err != nil {
		color.Error.Println("Error Make Folder: ", err)
	}
	n2, err := sscard.WriteBlockToFile(cardPhotoJpg, "./imgs/idcPhoto.jpg")
	if err != nil {
		color.Error.Println("Error WriteBlockToFile: ", err)
		return ReponseData{}, err
	}
	color.Info.Printf("wrote %d bytes\n", n2)

	return ReponseData{
		CID:        string(cidTrim),
		FullNameTh: fullnameTrim1,
		FullNameEn: string(fullnameEN),
		Birth:      string(birthTrim),
		Gender:     string(genderTrim),
		IssueExp:   string(issueExpTrim),
		Image:      c.BaseURL() + "/public/idcPhoto.jpg",
	}, nil
}

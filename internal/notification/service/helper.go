package service

import (
	"fmt"

	pb "github.com/wirsal/project-aegis/api/protos"
	"google.golang.org/protobuf/types/known/structpb"
)

/*TODO dynamic recipeint get from database */
func getRecipient(carnumber string) *Recipient {
	recipient := &Recipient{
		fcm_token:    "ch30ZNhq7kDgnPqChvHg6W:APA91bGn9oH5JuUh4BBV2gF_0B4dZXLfF2Yo94xFX7dr_w5awHjEZBfTpOjg1IwI1F-C6Ap9KB6lZGb2tFet4W3jEvviJ6q1aMzI9pNu_4iaqbpj12-FNfw",
		email:        "email@123.com",
		phone_number: "0812345678910",
	}

	return recipient
}

func getPayloadFCM(trxData *pb.Transaction) (*structpb.Struct, error) {

	payloadData := map[string]interface{}{
		"title": "Peringatan Keamanan Kartu Anda",
		"body":  fmt.Sprintf("Terdeteksi transaksi mencurigakan sebesar %.2f di %s.", trxData.GetTrxAmount(), trxData.GetTrxMerchantName()),
	}

	payloadStruct, err := structpb.NewStruct(payloadData)
	if err != nil {
		return nil, fmt.Errorf("failed to create payload struct: %w", err)
	}

	return payloadStruct, nil
}

func getPayloadEmail(trxData *pb.Transaction, template string) (*structpb.Struct, error) {
	return nil, nil
}

func getPayloadWa(trxData *pb.Transaction, template string) (*structpb.Struct, error) {
	return nil, nil
}

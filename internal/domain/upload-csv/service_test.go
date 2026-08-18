package uploadcsv

import "testing"

func TestMapHeaders(t *testing.T) {
	mappings, front, back, hasUID, err := mapHeaders([]string{
		"frn_name",
		"frn_img_logo",
		"frn_uid",
		"bck_address",
		"bck_uid",
		"trk1_data",
		"trk2_data",
		"trk3_data",
		"ignored",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(mappings) != 8 || len(front) != 3 || len(back) != 2 || !hasUID {
		t.Fatalf("unexpected mapping result: mappings=%d front=%d back=%d uid=%v", len(mappings), len(front), len(back), hasUID)
	}
	if mappings[1].contentType != "image" || !mappings[1].isImage {
		t.Fatalf("image column was not classified correctly: %+v", mappings[1])
	}
	if !mappings[2].isUID || mappings[2].uidBlock != -1 {
		t.Fatalf("front UID metadata is incorrect: %+v", mappings[2])
	}
	if mappings[5].trackNo != 1 || mappings[7].trackNo != 3 {
		t.Fatalf("track columns were not classified correctly")
	}
}

func TestMapHeadersRejectsDuplicateMappedColumns(t *testing.T) {
	_, _, _, _, err := mapHeaders([]string{"frn_name", "FRN_NAME"})
	if err == nil {
		t.Fatal("expected duplicate mapped column error")
	}
}

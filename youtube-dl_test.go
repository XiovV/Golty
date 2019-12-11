package main

//
//var targets1 = []DownloadTarget{
//	{
//		URL: "https://www.youtube.com/user/NewRetroWave/",
//	},
//	{
//		URL: "https://www.youtube.com/user/HungOverGargoyle",
//	},
//}
//
//func TestGetMetadata(t *testing.T) {
//	assert := assert.New(t)
//	for _, target := range targets1 {
//		targetMetadata, err := target.GetMetadata()
//		assert.Nil(err)
//		if target.URL == "https://www.youtube.com/user/NewRetroWave/" {
//			assert.Equal("NewRetroWave", targetMetadata.Uploader)
//		} else if target.URL == "https://www.youtube.com/user/HungOverGargoyle" {
//			assert.Equal("Electronic Gems", targetMetadata.Uploader)
//		}
//	}
//}
//
//func TestDownloadTarget_GetLatestVideo(t *testing.T) {
//	assert := assert.New(t)
//	for _, target := range targets1 {
//		latestVideoId, err := target.GetLatestVideo()
//		assert.Nil(err)
//		assert.Equal(11, len(latestVideoId), "length of latestVideoId should be 11 because that's how long video id's are on youtube")
//	}
//}

# golang-pipVideo-mixAudio

### build
docker build -t oceanwu/golang-pipvideo-mixaudio .

### Usage
docker run -d -p 8089:8089 -v $(pwd)/data:/go/src/github.com/mind-stack-cn/golang-pipVideo-mixAudio/data oceanwu/golang-pipvideo-mixaudio 

## request param
````
type PostParams struct {
	AudioUri0   string	// 待叠加音频0
	AudioUri1   string	// 待叠加音频1
	VideoUri0   string	// 待叠加视频0
	VideoUri1   string	// 待叠加视频1
	CallBackUrl string	// 回调请求url,POST
}
````

## response data
````
type MixedResult struct {
	MixedAudio interface{}
	MixedVideo interface{}
}
````

## test
curl -X POST -d '{"audioUri0": "https://raw.githubusercontent.com/mind-stack-cn/golang-pipVideo-mixAudio/master/test/testmedia/1.aac", "audioUri1": "https://raw.githubusercontent.com/mind-stack-cn/golang-pipVideo-mixAudio/master/test/testmedia/2.aac", "videoUri0": "https://raw.githubusercontent.com/mind-stack-cn/golang-pipVideo-mixAudio/master/test/testmedia/1.mp4", "videoUri1": "https://raw.githubusercontent.com/mind-stack-cn/golang-pipVideo-mixAudio/master/test/testmedia/2.mp4", "callBackUri": ""}' http://127.0.0.1:8089

````
{
    "MixedAudio":{
        "uri":"/155678c0/186d/4a5f/9a93/b3ab577f3351.aac","size":775739,"fileType":"audio","duration":47.778948
    },
    "MixedVideo":{
        "uri":"/6391b07a/38b6/41cf/b9d2/4bb6fb2898eb.mp4","size":946246,"fileType":"video","duration":10.032,
        "thumbnail":
            {
                "uri":"/6391b07a/38b6/41cf/b9d2/4bb6fb2898eb.jpg","size":15638,"fileType":"image","width":640,"height":360
            }
    }
}
````

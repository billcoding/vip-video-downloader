# vip-video-downloader
VIP Video Downloader, such as: iqiyi, youku, qq, ...etc.

# usage
- Download

vip-video-downloader download [flags] URL 

- Merge

vip-video-downloader merge [flags] DIRECTORY

# examples

## download

- normal: `vip-video-downloader download "https://abc.com/v/xyz.html"`
- output dir: `vip-video-downloader download -d "/to/path" "https://abc.com/v/xyz.html"`
- output file: `vip-video-downloader download -o "my-video" "https://abc.com/v/xyz.html"`
- no convert(.ts): `vip-video-downloader download -c=F "https://abc.com/v/xyz.html"`
- convert to mkv format: `vip-video-downloader download -f="mkv" "https://abc.com/v/xyz.html"`
- special FFmpeg path: `vip-video-downloader download -F="/to/path/ffmpeg" "https://abc.com/v/xyz.html"`
- m3u8 URL: `vip-video-downloader download -m "https://abc.com/index.m3u8"`
- m3u8 File: `vip-video-downloader download -M "/to/path/v.m3u8"`

## merge
- normal: `vip-video-downloader merge -o "output.tmp" /to/path`

# Copyright

__版权所有：2022 ©billcoding__

__注意： 本软件仅供学习使用，请勿用于非法用途，否则由使用者承担一切责任！推荐购买视频平台VIP包月观看会员视频！__
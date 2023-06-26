package utils

import "testing"

func TestFileMd5(t *testing.T) {
	t.Log(FileMd5("/home/hq/project/myproject/kuajing_work/kuajing/public/upload/upload/127669483824231005.jpg"))
	t.Log(FileMd5("/home/hq/project/myproject/kuajing_work/kuajing/public/upload/upload/127810624871739997.jpg"))
	t.Log(FileMd5("/home/hq/project/myproject/kuajing_work/kuajing/public/upload/upload/127812860871979613.jpg"))
	t.Log(FileMd5("/home/hq/project/myproject/kuajing_work/kuajing/public/upload/upload/130416611680072285.jpg"))
	t.Log(FileMd5("/home/hq/rumenz.img"))
}

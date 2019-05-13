package blurhash_test

var testFixtures = []struct {
	file         string
	hash         string
	xComp, yComp int
}{
	{"fixtures/test.png", "LFE.@D9F01_2%L%MIVD*9Goe-;WB", 4, 3},
	{"fixtures/octocat.png", "LNAdApj[00aymkj[TKay9}ay-Sj[", 4, 3},
	{"", "LNMF%n00%#MwS|WCWEM{R*bbWBbH", 4, 3},
	{"", "KJG8_@Dgx]_4V?xuyE%NRj", 3, 3},
}

package channel

import (
	"fmt"
	"net/url"
	"testing"
)

func TestURLParse(t *testing.T) {
	rawURL := "https://sluiceyf.titan.mgtv.com/c1/2021/11/24_0/70D25AF0FCF6481864ED3910FB47CBDF_20211124_1_1_605_mp4/5F480BC2C8CACE0EC124CE570F5AFEDF.m3u8?arange=0&jxname=jianghu&pm=4A_z~cyrTMY7p00xpSyNR47EcR5wpcIpKDcVSOYh3zgz0avZkcWsWUT4EBN9LwqdkuwfTfTn5xAoW9Gbf2qNKttuCChOrTmG7Ehmlv8rz1g3xZhzgXukgTjES87tXRnMFAYSjov1uRGZsApB1KJluoFg3AJpGwxqReYIDlJ16DQ4bP1~4~mroNCLlGxrG~qGIC2YEAIH91dkrX0vgGr36yxVtlOp~licOr1OtFX2mR9QqAtAnYag4Ny2EmQRj7LgdFkjIM5MSYtBVYUIUjOFfv1NR4hvcfvMCN_si3q7u~g8ohUmBpIP6VD8cLBzD5ImkejXjCeJVimTjVOOfw9Vg2WzI5maQWN2BdPKBOhgSsxBdOnDbnIZ6_O1TmSNyRcUl~5KoBMoNXkBCSJZy6H0JPSozYJBukLARSCFo8SOu3IDLnP7aGKrPDmsW1szD4h7~ihCCYjKb~wtfc_VUHhGf9gcRxQ-&uid=9681952cec131090e092c7ac077568f4&vcdn=0&scid=25117&ruid=780449c7fbd34660&cp=hpmMlJaPwss~"
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	fmt.Println(parsedURL.Path)
	fmt.Println(parsedURL.Query())
}

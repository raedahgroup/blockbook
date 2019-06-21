package dcr

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"
)

var (
	testTx1, testTx2 bchain.Tx

	testTxPacked1 = "00003dcb8bb8bf943c00000000007b22686578223a2230313030303030303031323337323536386665383064326639623261623137323236313538646435373332643939323664633730353337316561663430616237343863396533643937323032303030303030303166666666666666663032363434623235326430303030303030303030303031393736613931346138363266383337333363633336386633383661363531653033643834346135626436313136643538386163616364663633303930303030303030303030303031393736613931343931646335643138333730393339623334313436303361303732396263623361333865346566373638386163303030303030303030303030303030303031653438643839333630303030303030306262336430303030303230303030303036613437333034343032323033373865313434326363313766613765343931383435313837313365656464333065313365343231343765303737383539353537646136666662626434306337303232303566383535363363323862363238376639633931313065363836346464313861636664393264383535303965613834363931336332386236653861376639343030313231303262626264376161646566333366326432626464396230633562613237383831356635643636613661303164326330313966623733663639373636323033386235222c2274786964223a2231333261636235623437346234356238333066373936316339316338376535336363653361333761366336663062303933336363646630333935633831613661222c2276657273696f6e223a312c226c6f636b74696d65223a302c2276696e223a5b7b22636f696e62617365223a22222c2274786964223a2237326439653363393438623730616634656137313533373064633236393932643733643538643135323637326231326139623266306465383866353637323233222c22766f7574223a322c22736372697074536967223a7b22686578223a22227d2c2273657175656e6365223a343239343936373239352c22616464726573736573223a5b5d7d5d2c22766f7574223a5b7b2256616c7565536174223a3735373431383835322c2276616c7565223a302c226e223a302c227363726970745075624b6579223a7b22686578223a223736613931346138363266383337333363633336386633383661363531653033643834346135626436313136643538386163222c22616464726573736573223a5b225473674e555a4b456e55684641534c45536a37665652546b677565335152395441655a225d7d7d2c7b2256616c7565536174223a3135373534303236382c2276616c7565223a302c226e223a312c227363726970745075624b6579223a7b22686578223a223736613931343931646335643138333730393339623334313436303361303732396263623361333865346566373638386163222c22616464726573736573223a5b225473654b4e53575962417a61476f67706e4e6e32357465547a353350546b3373675075225d7d7d5d2c22626c6f636b74696d65223a313533353633323637307d"
	testTxPacked2 = "00003df38bb8bfec6c00000000007b22686578223a2230313030303030303031633536643830373536656161376663366533353432623239663539366336306139626363393539636630346435663665366231323734396532343165636532393032303030303030303166666666666666663032636632306234326430303030303030303030303031393736613931343037393964616133636433366234346465663232303838363830326562393965313063346137633438386163306332356337303730303030303030303030303031393736613931343062313032646562333331343231333136346362363332323231313232353336353635383430376538386163303030303030303030303030303030303031616661383762333530303030303030306533336430303030303030303030303036613437333034343032323031666633343265356161353562363033303137316638353732393232316361306238313933383832366363303934343962373737353265366533623631356265303232303238316531363062363138653537333236623935613065306333616337613531336264303431616261363363626163653266373139313965313131636664626130313231303239306138646536363635633863616163326262386361316161626433646330396133333466393937663937626438393437373262316535316361623030336439222c2274786964223a2263616633346339333464346333366234313063303236353232326230363966353265326466343539656262303964363739376136333563656565306564643630222c2276657273696f6e223a312c226c6f636b74696d65223a302c2276696e223a5b7b22636f696e62617365223a22222c2274786964223a2232396365316532343965373431323662366535663464663039633935636339623061633639366635323932623534653363363766616136653735383036646335222c22766f7574223a322c22736372697074536967223a7b22686578223a22227d2c2273657175656e6365223a343239343936373239352c22616464726573736573223a5b5d7d5d2c22766f7574223a5b7b2256616c7565536174223a3736363737393539392c2276616c7565223a302c226e223a302c227363726970745075624b6579223a7b22686578223a223736613931343037393964616133636433366234346465663232303838363830326562393965313063346137633438386163222c22616464726573736573223a5b22547352694b577353397563617159447739716867364e756b54746853354c7754526e76225d7d7d2c7b2256616c7565536174223a31333034393136362c2276616c7565223a302c226e223a312c227363726970745075624b6579223a7b22686578223a223736613931343062313032646562333331343231333136346362363332323231313232353336353635383430376538386163222c22616464726573736573223a5b2254735332644871455359317666666a6464706f31564d5462774c6e44737066456a3557225d7d7d5d2c22626c6f636b74696d65223a313533353633383332367d"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:       "01000000012372568fe80d2f9b2ab17226158dd5732d9926dc705371eaf40ab748c9e3d9720200000001ffffffff02644b252d0000000000001976a914a862f83733cc368f386a651e03d844a5bd6116d588acacdf63090000000000001976a91491dc5d18370939b3414603a0729bcb3a38e4ef7688ac000000000000000001e48d893600000000bb3d0000020000006a4730440220378e1442cc17fa7e49184518713eedd30e13e42147e077859557da6ffbbd40c702205f85563c28b6287f9c9110e6864dd18acfd92d85509ea846913c28b6e8a7f940012102bbbd7aadef33f2d2bdd9b0c5ba278815f5d66a6a01d2c019fb73f697662038b5",
		Blocktime: 1535632670,
		Txid:      "132acb5b474b45b830f7961c91c87e53cce3a37a6c6f0b0933ccdf0395c81a6a",
		LockTime:  0,
		Version:   1,
		Vin: []bchain.Vin{
			{
				Txid:      "72d9e3c948b70af4ea715370dc26992d73d58d152672b12a9b2f0de88f567223",
				Vout:      2,
				Sequence:  4294967295,
				Addresses: []string{},
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat:  *big.NewInt(757418852),
				N:         0,
				JsonValue: json.Number("0"),
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a914a862f83733cc368f386a651e03d844a5bd6116d588ac",
					Addresses: []string{
						"TsgNUZKEnUhFASLESj7fVRTkgue3QR9TAeZ",
					},
				},
			},
			{
				ValueSat:  *big.NewInt(157540268),
				N:         1,
				JsonValue: json.Number("0"),
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a91491dc5d18370939b3414603a0729bcb3a38e4ef7688ac",
					Addresses: []string{
						"TseKNSWYbAzaGogpnNn25teTz53PTk3sgPu",
					},
				},
			},
		},
	}

	testTx2 = bchain.Tx{
		Hex:       "0100000001c56d80756eaa7fc6e3542b29f596c60a9bcc959cf04d5f6e6b12749e241ece290200000001ffffffff02cf20b42d0000000000001976a9140799daa3cd36b44def220886802eb99e10c4a7c488ac0c25c7070000000000001976a9140b102deb3314213164cb6322211225365658407e88ac000000000000000001afa87b3500000000e33d0000000000006a47304402201ff342e5aa55b6030171f85729221ca0b81938826cc09449b77752e6e3b615be0220281e160b618e57326b95a0e0c3ac7a513bd041aba63cbace2f71919e111cfdba01210290a8de6665c8caac2bb8ca1aabd3dc09a334f997f97bd894772b1e51cab003d9",
		Blocktime: 1535638326,
		Txid:      "caf34c934d4c36b410c0265222b069f52e2df459ebb09d6797a635ceee0edd60",
		LockTime:  0,
		Version:   1,
		Vin: []bchain.Vin{
			{
				Txid:      "29ce1e249e74126b6e5f4df09c95cc9b0ac696f5292b54e3c67faa6e75806dc5",
				Vout:      2,
				Sequence:  4294967295,
				Addresses: []string{},
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat:  *big.NewInt(766779599),
				N:         0,
				JsonValue: json.Number("0"),
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9140799daa3cd36b44def220886802eb99e10c4a7c488ac",
					Addresses: []string{
						"TsRiKWsS9ucaqYDw9qhg6NukTthS5LwTRnv",
					},
				},
			},
			{
				ValueSat:  *big.NewInt(13049166),
				N:         1,
				JsonValue: json.Number("0"),
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9140b102deb3314213164cb6322211225365658407e88ac",
					Addresses: []string{
						"TsS2dHqESY1vffjddpo1VMTbwLnDspfEj5W",
					},
				},
			},
		},
	}
}

func TestGetAddrDescFromAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH",
			args:    args{address: "TcrypGAcGCRVXrES7hWqVZb5oLJKCZEtoL1"},
			want:    "5463727970474163474352565872455337685771565a62356f4c4a4b435a45746f4c31",
			wantErr: false,
		},
		{
			name:    "P2PKH",
			args:    args{address: "TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd"},
			want:    "547366444c72526b6b3963695575776670326238506177776e756b59443779416a4764",
			wantErr: false,
		},
		{
			name:    "P2PKH",
			args:    args{address: "TsTevp3WYTiV3X1qjvZqa7nutuTqt5VNeoU"},
			want:    "547354657670335759546956335831716a765a7161376e75747554717435564e656f55",
			wantErr: false,
		},
	}
	parser := NewDecredParser(GetChainParams("testnet3"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestGetAddrDescFromVout(t *testing.T) {
	type args struct {
		vout bchain.Vout
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PK",
			args:    args{vout: bchain.Vout{ScriptPubKey: bchain.ScriptPubKey{Hex: "76a914936f3a56a2dd0fb3bfde6bc820d4643e1701542a88ac"}}},
			want:    "54736554683431516f356b594c3337614c474d535167346e67636f71396a7a44583659",
			wantErr: false,
		},
		{
			name:    "P2PK",
			args:    args{vout: bchain.Vout{ScriptPubKey: bchain.ScriptPubKey{Hex: "76a9144b31f712b03837b1303cddcb1ae9abd98da44f1088ac"}}},
			want:    "547358736a3161747744736455746e354455576b666f6d5a586e4a6151467862395139",
			wantErr: false,
		},
		{
			name:    "P2PK",
			args:    args{vout: bchain.Vout{ScriptPubKey: bchain.ScriptPubKey{Hex: "76a9140d85a1d3f77383eb3dacfd83c46e2c7915aba91d88ac"}}},
			want:    "54735346644c79657942776e68486978737367784b34546f4664763876525931793871",
			wantErr: false,
		},
	}
	parser := NewDecredParser(GetChainParams("testnet3"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromVout(&tt.args.vout)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromVout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromVout() = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestGetAddressesFromAddrDesc(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2PKH",
			args:    args{script: "5463727970474163474352565872455337685771565a62356f4c4a4b435a45746f4c31"},
			want:    []string{"TcrypGAcGCRVXrES7hWqVZb5oLJKCZEtoL1"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2PKH",
			args:    args{script: "547366444c72526b6b3963695575776670326238506177776e756b59443779416a4764"},
			want:    []string{"TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "P2PKH",
			args:    args{script: "547354657670335759546956335831716a765a7161376e75747554717435564e656f55"},
			want:    []string{"TsTevp3WYTiV3X1qjvZqa7nutuTqt5VNeoU"},
			want2:   true,
			wantErr: false,
		},
	}

	parser := NewDecredParser(GetChainParams("testnet3"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestPackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *DecredParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "pack-tx-1",
			args: args{
				tx:        testTx1,
				height:    15819,
				blockTime: 1535632670,
				parser:    NewDecredParser(GetChainParams("testnet3"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
		{
			name: "pack-tx-2",
			args: args{
				tx:        testTx2,
				height:    15859,
				blockTime: 1535638326,
				parser:    NewDecredParser(GetChainParams("testnet3"), &btc.Configuration{}),
			},
			want:    testTxPacked2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestUnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *DecredParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "unpack-tx-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewDecredParser(GetChainParams("testnet3"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   15819,
			wantErr: false,
		},
		{
			name: "unpack-tx-2",
			args: args{
				packedTx: testTxPacked2,
				parser:   NewDecredParser(GetChainParams("testnet3"), &btc.Configuration{}),
			},
			want:    &testTx2,
			want1:   15859,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %+v, want %+v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

package main

import "testing"

func Test_jsonEqual(t *testing.T) {
	type args struct {
		a            string
		b            string
		ignoreBodies map[string]struct{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "",
			args: args{
				a: `{"a":1,"b":2}`,
				b: `{"b":2,"a":1}`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"a":1,"b":2}`,
				b: `{"b":2,"a":2}`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"a":1,"b":2}`,
				b: `{"b":2,"a":1,"c":3}`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"a":1,"b":2}`,
				b: `{"b":2,"a":1,"c":3}`,
				ignoreBodies: map[string]struct{}{
					"c": {},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"language":"go"}
`,
				b: `{"language":"go"}
`,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"viewerId":"23926361","sessionId":"d7921fd3-a481-4107-8a04-c1b4c7b7470b","updatedResources":{"now":1661525999,"user":{"id":10208648,"isuCoin":7746,"lastGetRewardAt":1568271756,"lastActivatedAt":1661525999,"registeredAt":1567746156,"createdAt":1567746156,"updatedAt":1661525999},"userLoginBonuses":[{"id":47040139691,"userId":10208648,"loginBonusId":1,"lastRewardSequence":3,"loopCount":13,"createdAt":1567746156,"updatedAt":1661525999}],"userPresents":[{"id":1137736622062702592,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622062702594,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１３ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896896,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１４ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896898,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１５ヶ月突破プレゼントです１","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896900,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１６ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896902,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１７ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896904,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１８ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896906,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１９ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896908,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２０ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896910,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２１ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896912,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２２ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896914,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２３ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896916,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896918,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２.５周年プレゼントです！","createdAt":1661525999,"updatedAt":1661525999}]}}`,
				b: `{"viewerId":"23926361","sessionId":"6a7b5da6-5bf0-40bf-86d1-6f827aa26611","updatedResources":{"now":1661525999,"user":{"id":10208648,"isuCoin":7746,"lastGetRewardAt":1568271756,"lastActivatedAt":1661525999,"registeredAt":1567746156,"createdAt":1567746156,"updatedAt":1661525999},"userLoginBonuses":[{"id":47040139691,"userId":10208648,"loginBonusId":1,"lastRewardSequence":3,"loopCount":13,"createdAt":1567746156,"updatedAt":1661525999}],"userPresents":[{"id":1137871631067123712,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123714,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１３ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123716,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１４ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123718,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１５ヶ月突破プレゼントです１","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123720,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１６ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123722,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１７ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318016,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１８ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318018,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１９ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318020,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２０ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318022,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２１ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318024,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２２ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318026,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２３ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318028,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318030,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２.５周年プレゼントです！","createdAt":1661525999,"updatedAt":1661525999}]}}`,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"viewerId":"23926361","sessionId":"d7921fd3-a481-4107-8a04-c1b4c7b7470b","updatedResources":{"now":1661525999,"user":{"id":10208648,"isuCoin":7746,"lastGetRewardAt":1568271756,"lastActivatedAt":1661525999,"registeredAt":1567746156,"createdAt":1567746156,"updatedAt":1661525999},"userLoginBonuses":[{"id":47040139691,"userId":10208648,"loginBonusId":1,"lastRewardSequence":3,"loopCount":13,"createdAt":1567746156,"updatedAt":1661525999}],"userPresents":[{"id":1137736622062702592,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622062702594,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１３ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896896,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１４ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896898,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１５ヶ月突破プレゼントです１","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896900,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１６ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896902,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１７ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896904,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１８ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896906,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１９ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896908,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２０ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896910,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２１ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896912,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２２ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896914,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２３ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896916,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896918,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２.５周年プレゼントです！","createdAt":1661525999,"updatedAt":1661525999}]}}`,
				b: `{"viewerId":"23926361","sessionId":"6a7b5da6-5bf0-40bf-86d1-6f827aa26611","updatedResources":{"now":1661525999,"user":{"id":10208648,"isuCoin":7746,"lastGetRewardAt":1568271756,"lastActivatedAt":1661525999,"registeredAt":1567746156,"createdAt":1567746156,"updatedAt":1661525999},"userLoginBonuses":[{"id":47040139691,"userId":10208648,"loginBonusId":1,"lastRewardSequence":3,"loopCount":13,"createdAt":1567746156,"updatedAt":1661525999}],"userPresents":[{"id":1137871631067123712,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123714,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１３ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123716,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１４ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123718,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１５ヶ月突破プレゼントです１","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123720,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１６ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123722,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１７ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318016,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１８ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318018,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１９ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318020,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２０ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318022,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２１ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318024,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２２ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318026,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２３ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318028,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318030,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２.５周年プレゼントです！","createdAt":1661525999,"updatedAt":1661525999}]}}`,
				ignoreBodies: map[string]struct{}{
					"sessionId": {},
					"id":        {},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"viewerId":"23926361","sessionId":"d7921fd3-a481-4107-8a04-c1b4c7b7470b","updatedResources":{"now":1661525999,"user":{"id":10208648,"isuCoin":7746,"lastGetRewardAt":1568271756,"lastActivatedAt":1661525999,"registeredAt":1567746156,"createdAt":1567746156,"updatedAt":1661525999},"userLoginBonuses":[{"id":47040139691,"userId":10208648,"loginBonusId":1,"lastRewardSequence":3,"loopCount":13,"createdAt":1567746156,"updatedAt":1661525999}],"userPresents":[{"id":1137736622062702592,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622062702594,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１３ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896896,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１４ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896898,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１５ヶ月突破プレゼントです１","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896900,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１６ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896902,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１７ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896904,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１８ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896906,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１９ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896908,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２０ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896910,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２１ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896912,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２２ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896914,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２３ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896916,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896918,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２.５周年プレゼントです！","createdAt":1661525999,"updatedAt":1661525999}]}}`,
				b: `{"viewerId":"23926361","sessionId":"6a7b5da6-5bf0-40bf-86d1-6f827aa26611","updatedResources":{"now":1661525999,"user":{"id":10208648,"isuCoin":7746,"lastGetRewardAt":1568271756,"lastActivatedAt":1661525999,"registeredAt":1567746156,"createdAt":1567746156,"updatedAt":1661525999},"userLoginBonuses":[{"id":47040139691,"userId":10208648,"loginBonusId":1,"lastRewardSequence":3,"loopCount":13,"createdAt":1567746156,"updatedAt":1661525999}],"userPresents":[{"id":1137871631067123712,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123714,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１３ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123716,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１４ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123718,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１５ヶ月突破プレゼントです１","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123720,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１６ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631067123722,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１７ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318016,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１８ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318018,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１９ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318020,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２０ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318022,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２１ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318024,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２２ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318026,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２３ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318028,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137871631071318030,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２.５周年プレゼントです！","createdAt":1661525999,"updatedAt":1661525999}]}}`,
				ignoreBodies: map[string]struct{}{
					"sessionId":                        {},
					"updatedResources.userPresents.id": {},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"viewerId":"23926361","sessionId":"d7921fd3-a481-4107-8a04-c1b4c7b7470b","updatedResources":{"now":1661525999,"user":{"id":10208648,"isuCoin":7746,"lastGetRewardAt":1568271756,"lastActivatedAt":1661525999,"registeredAt":1567746156,"createdAt":1567746156,"updatedAt":1661525999},"userLoginBonuses":[{"id":47040139691,"userId":10208648,"loginBonusId":1,"lastRewardSequence":3,"loopCount":13,"createdAt":1567746156,"updatedAt":1661525999}],"userPresents":[{"id":1137736622062702592,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622062702594,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１３ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896896,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１４ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896898,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１５ヶ月突破プレゼントです１","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896900,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１６ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896902,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１７ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896904,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１８ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896906,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１９ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896908,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２０ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896910,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２１ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896912,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２２ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896914,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２３ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896916,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137736622066896918,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２.５周年プレゼントです！","createdAt":1661525999,"updatedAt":1661525999}]}}`,
				b: `{"viewerId":"23926361","sessionId":"9a19c55a-b732-4c3e-88fa-5b3b495c1e88","updatedResources":{"now":1661525999,"user":{"id":10208648,"isuCoin":7746,"lastGetRewardAt":1568271756,"lastActivatedAt":1661525999,"registeredAt":1567746156,"createdAt":1567746156,"updatedAt":1661525999},"userLoginBonuses":[{"id":47040139691,"userId":10208648,"loginBonusId":1,"lastRewardSequence":3,"loopCount":13,"createdAt":1567746156,"updatedAt":1661525999}],"userPresents":[{"id":1137881698990559232,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559234,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１３ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559236,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１４ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559238,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１５ヶ月突破プレゼントです１","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559240,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１６ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559242,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１７ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559244,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１８ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559246,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"１９ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559248,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２０ヶ月突破プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559250,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２１ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559252,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２２ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559254,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":2000,"presentMessage":"２３ヶ月プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698990559256,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２周年記念プレゼントです！","createdAt":1661525999,"updatedAt":1661525999},{"id":1137881698994753536,"userId":10208648,"sentAt":1661525999,"itemType":1,"itemId":1,"amount":6000,"presentMessage":"２.５周年プレゼントです！","createdAt":1661525999,"updatedAt":1661525999}]}}`,
				ignoreBodies: map[string]struct{}{
					"sessionId":                        {},
					"updatedResources.userPresents.id": {},
					"session.id":                       {},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				a: `{"session":{"id":1137736635526418432,"userId":123456,"viewerId":"","sessionId":"1648e138-f4ab-4741-b78f-f7af549d7179","expiredAt":1691414351,"createdAt":1691327951,"updatedAt":1691327951}}`,
				b: `{"session":{"id":1137894998860107776,"userId":123456,"viewerId":"","sessionId":"d62f4600-9148-4dec-970a-ac4545e52788","expiredAt":1691452108,"createdAt":1691365708,"updatedAt":1691365708}}`,
				ignoreBodies: map[string]struct{}{
					"session.id":        {},
					"session.sessionId": {},
					"session.expiredAt": {},
					"session.createdAt": {},
					"session.updatedAt": {},
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonEqual(tt.args.a, tt.args.b, tt.args.ignoreBodies)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonEqual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("jsonEqual() got = %v, want %v", got, tt.want)
			}
		})
	}
}

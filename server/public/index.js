const ajax = async (url, method, params = null) => {
    let obj = { method, headers: { "Content-Type": "application/json; charset=utf-8" } };
    if (params != null) obj["params"] = params;
    let resp = await fetch(url, obj);
    let json = await resp.json();
    return json;
};
const { createApp } = Vue;
createApp({
    data() {
        return {
            info: `ChallengeApp by j.rebhan@gmail.com (${new Date().toLocaleDateString()} ${new Date().toLocaleTimeString()})`,
            result: "",

            url: '/api/v1/dbversion',
            endpoints: [
                { text: 'testcase 0: about', value: '/api/v1/about' },
                { text: 'testcase 1: db version', value: '/api/v1/dbversion' },
                { text: 'testcase 2: top_zones by pickups', value: '/api/v1/top-zones/pickups' },
                { text: 'testcase 3: top_zones by dropoffs', value: '/api/v1/top-zones/dropoffs' },
                { text: 'testcase 4: error message', value: '/api/v1/errorexample' },
                { text: 'testcase 5: error message top-zones', value: '/api/v1/top-zones/errorexamle' },
                { text: 'testcase 6: zone-trips by zone and date', value: '/api/v1/zone-trips/7/2018-01-06' },
                { text: 'testcase 7: zone-trips error test: no params', value: '/api/v1/zone-trips/' },
                { text: 'testcase 8: zone-trips error test: wrong params', value: '/api/v1/zone-trips/errortest' },
                { text: 'testcase 9: list-yellow limit 100 ', value: '/api/v1/list-yellow/limit=100' },
                { text: 'testcase 10: list-yellow limit 10, offeset 10', value: '/api/v1/list-yellow/limit=10/offset=10' },
                { text: 'testcase 11: list-yellow, test no params behaviour (error expected)', value: '/api/v1/list-yellow' },
                { text: 'testcase 12: list_yellow sort by pu_location asc', value: '/api/v1/list-yellow/sort=pu_datetime.asc/sort=pu_locationid.asc' },
                { text: 'testcase 13: list_yellow sort by pu_location desc', value: '/api/v1/list-yellow/sort=pu_datetime.asc/sort=pu_locationid.desc' },
                { text: 'testcase 14: list_yellow sort by pu_datetime > 2018-01-20, sort  pu_locationid asc, filter do_locationid >= 100, offset = 30', value: '/api/v1/list-yellow/sort=pu_datetime.asc/limit=10/filter=pu_datetime:gt:2018-01-20/sort=pu_locationid.asc/filter=do_locationid:gte:100/offset=30' },
                { text: 'testcase 15: list_yellow changed parameter order: offset = 10, sort pu_datetime desc, sort pu_locationid desc', value: '/api/v1/list-yellow/limit=10/offset=10/sort=pu_datetime.desc/sort=pu_locationid.desc' },
            ]
        };
    },
    methods: {
        async api(url, type) {
            let data = await ajax(url, type);
            this.result = JSON.stringify(data, null, " ");
        },
        clear() {
            this.result = "";
            this.url = '/api/v1/dbversion';
        }
    },
    mounted() {
        console.log("ChallengeApp Testbed started");
    },
}).mount("#challengeapp");
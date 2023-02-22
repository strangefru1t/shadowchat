package templates

var DashboardHTML string = `<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8" />
    <title>Dashboard</title>
    <link rel="stylesheet" href="/style.css">
</head>

<style>

</style>

<body style="max-width:500px ;margin:0 auto">
    <h2 style="text-align: center">DASHBOARD</h2>
    <button class="tab" onmousedown="openCity(event,'Chats')">Chats</button>
    <button class="tab" onmousedown="openCity(event,'Overlays')">Overlays</button>
    <button class="tab" onmousedown="openCity(event,'Settings')">Settings</button>
  
  <div id="Chats" class="panel">
    <script>
        var dateformat = { month: 'short', day: 'numeric', hour: 'numeric', minute: 'numeric' };
        fetch('/api/chatlog')
            .then(function (response) {
                return response.json();
            })
            .then(function (data) {
                appendData(data);
            })
            .catch(function (err) {
                console.log('error: ' + err);
            });
        function appendData(data) {
            data.sort((a,b) => a.timestamp > b.timestamp ? -1 : 1)
            var mainContainer = document.getElementById("Chats");
            for (var i = 0; i < data.length; i++) {
                var div = document.createElement("blockquote");
                div.classList.add('card');
                var time = new Date(data[i].timestamp).toLocaleTimeString("en-US",dateformat);
                div.innerHTML = ` + "`" + `<h3><span style="color:${id2color(data[i].id)}">${data[i].name}</span> <span class=grey>sent</span> ${data[i].amount} <span class=grey>(${time})</span></h3><p>${data[i].message}</p>` + "`" + `;
                mainContainer.appendChild(div);
            }
        }
        function id2color(stringInput) {
            let stringUniqueHash = [...stringInput].reduce((acc, char) => {
                return char.charCodeAt(0) + ((acc << 5) - acc);
            }, 0);
            return ` + "`" + `hsl(${stringUniqueHash % 360}, 100%, 60%)` + "`" + `;
        }
    </script>
  </div>
  <div id="Control Panel" class="panel" style="display:none">
    <h2>Control Panel</h2>
  </div>

  <div id="Settings" class="panel" style="display:none">
<form onsubmit="return post()">
  <label for="mindono">Minimum Donation (XMR):</label>
  <input type="number" id="mindono" name="mindono" step="0.001">
  <label for="maxchars">Max Donation Characters:</label>
  <input type="number" id="maxchars" name="maxchars">
  <label for="donogoal">Donation Goal (XMR):</label>
  <input type="number" id="donogoal" name="donogoal" step="0.001">
  <label for="donogoalhisthours">Donation Goal History (hours):</label>
  <input type="number" id="donogoalhisthours" name="donogoalhisthours">
  <label for="sltoken">Streamlabs Token:</label>
  <input type="text" id="streamlabstoken" name="streamlabstoken">
  <label for="alertinusd">Streamlabs Alert in USD:</label>
  <input type="checkbox" id="alertinusd" value=false name="alertinusd">
<input type="submit" id="save" value="Save" /></form>
</form> 
    <script>
        fetch('/api/settings')
            .then(function (response) {
                return response.json();
            })
            .then(function (data) {
				document.getElementById("mindono").value = data.mindono;	
				document.getElementById("maxchars").value = data.maxchars;	
				document.getElementById("streamlabstoken").value = data.streamlabstoken;	
				document.getElementById("alertinusd").checked = data.alertinusd;	
				document.getElementById("donogoal").value = data.donogoal;	
				document.getElementById("donogoalhisthours").value = data.donogoalhisthours;	
                console.log(data);
            })
            .catch(function (err) {
                console.log('error: ' + err);
            });
function post() {
    var data = new FormData();
    var checkedTrue = document.getElementById("alertinusd").checked === true;
    var formJson = JSON.stringify({
        "maxchars": parseInt(document.getElementById("maxchars").value),
        "alertinusd": checkedTrue,
        "streamlabstoken": document.getElementById("streamlabstoken").value,
        "mindono": parseFloat(document.getElementById("mindono").value),
        "donogoalhisthours": parseFloat(document.getElementById("donogoalhisthours").value),
        "donogoal": parseFloat(document.getElementById("donogoal").value)
    });
    fetch("/api/settings", {
        method: "POST",
        body: formJson
    }).then((result) => {
        if (result.status != 200) {
            throw new Error("Bad Server Response")
        }
        return result.json()
    }).then((response) => {
        var resp = JSON.parse(response);
        document.getElementById('save').value = resp.message;
    }).catch((error) => {
        console.log(error)
    });
    return false
}

    </script>
  </div>

  <div id="Overlays" class="panel" style="display:none">
          <h2>Overlays (WIP: more to come)</h2>
    <a href=/overlay/goal>/overlay/goal</a>
    <p>paste link into obs, interact with it to login, then refresh. OBS should hold onto the auth cookie</p>
  </div>

<script>
function openCity(evt, tab) {
  var i, x, tablinks;
  x = document.getElementsByClassName("panel");
  for (i = 0; i < x.length; i++) {
    x[i].style.display = "none";
  }
  tablinks = document.getElementsByClassName("tab");
  document.getElementById(tab).style.display = "block";
}
</script>
</body>

</html>`

var LoginHTML string = `<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8" />
    <title>login</title>
    <link rel="stylesheet" href="/style.css">
</head>

<style>

</style>

<body style="max-width:600px ;margin:0 auto">
    <h2 style="text-align: center">LOGIN</h2>
    <form method="post">
        <label>Username: </label>
        <input type="text" placeholder="Enter Username" name="username" required>
        <label>Password: </label>
        <input type="password" placeholder="Enter Password" name="password" required>
        <button type="submit">Login</button>
    </form>
</body>

</html>`

var RPCCONF string = `daemon-address=https://xmr-node.cakewallet.com:18081
wallet-file=/opt/shadowchat/wallet
rpc-bind-port=18082`

var StyleCSS string = `.waiting {
    display: none;
    position: relative;
    width: 80px;
    margin: 0 auto;
    height: 80px;
}

.waiting div {
    position: absolute;
    border: 5px solid #70799b;
    opacity: 1;
    border-radius: 50%;
    animation: waiting 1s cubic-bezier(0, 0.2, 0.8, 1) infinite;
}

.waiting div:nth-child(2) {
    animation-delay: -0.5s;
}

@keyframes waiting {
    0% {
        top: 36px;
        left: 36px;
        width: 0;
        height: 0;
        opacity: 0;
    }

    4.9% {
        top: 36px;
        left: 36px;
        width: 0;
        height: 0;
        opacity: 0;
    }

    5% {
        top: 36px;
        left: 36px;
        width: 0;
        height: 0;
        opacity: 1;
    }

    100% {
        top: 0px;
        left: 0px;
        width: 72px;
        height: 72px;
        opacity: 0;
    }
}

body {
    background: #12151C;
    color: ghostwhite;
    font-family: Sans-Serif;
    overflow: -moz-scrollbars-none;
}

body::-webkit-scrollbar {
    width: 0 !important
}

small {
    color: #70799b88;
    display: block;
    text-align: center
}

small a {
    color: #70799b88
}

a {
    color: #70799bff;
    text-decoration: none;
    outline: none
}

.form,.centered {
    max-width: 360px;
    margin: 0px auto;
}

::placeholder {
    color: #70799b;
}

label {
    -moz-user-select: none;
    -khtml-user-select: none;
    -webkit-user-select: none;
}

:focus {
    border: 1px solid #70799bdd;
}

textarea {
    resize: vertical;
    overflow: hidden;
    height: 64px;
    max-height: 128px;
}

.copy {
    user-select: all
}
.green { color: #6adda5 }
.blue { color: #96b4f1 }
.grey { color: #70799b }
input,textarea,blockquote,img,button,form a {
    display: block;
    word-break: break-all;
    outline: none;
    box-sizing: border-box;
    font-family: Sans-Serif;
    background: #262b3b;
    border: 1px solid #70799b60;
    color: white;
    min-height: 36px;
    padding: 8px;
    font-size: 16px;
    width: calc(100% - 10px);
    margin: 5px;
    border-radius: 4px;
}
input[type='checkbox'] { 
    display: inline;
    width: 16px;
    vertical-align: middle;
}
.card small, .card h3 {
  display: inline;
}
.card p {
   margin: 0px 0px 0px 0px;
   padding: 8px;
}
button {
  display: inline;
  width: auto;
}


form a {
    padding: 0px 4px 0px 4px;
    display: inline;
}`

var SCAPIJS string = `fetch('limits').then(function(response){return response.json()}).then(function(json){var name=document.getElementById("name");var amount=document.getElementById("amount");var amountlabel=document.getElementById("amountlabel");var message=document.getElementById("message");name.value="";amount.value="";message.value="";amount.min=json.MinimumDonationXMR;amount.placeholder=json.MinimumDonationXMR+" Minimum ($"+(json.XMRUSD*amount.min).toFixed(2)+")";message.placeholder=json.MaximumMessageChars+" Characters";message.maxLength=json.MaximumMessageChars;document.getElementById('three').addEventListener('click',function(e){amount.value=(3/json.XMRUSD).toFixed(4)});document.getElementById('five').addEventListener('click',function(e){amount.value=(5/json.XMRUSD).toFixed(4)});document.getElementById('ten').addEventListener('click',function(e){amount.value=(10/json.XMRUSD).toFixed(4)})}).catch(function(err){console.log('error: '+err)});function send(){document.getElementById("pay").remove();document.getElementById("three").remove();document.getElementById("five").remove();document.getElementById("ten").remove();var data=new FormData();var formJson=JSON.stringify({"name":document.getElementById("name").value,"message":document.getElementById("message").value,"amount":parseFloat(document.getElementById("amount").value)});document.getElementById("waiting").style.display="block";fetch("pay",{method:"POST",body:formJson}).then((result)=>{if(result.status!=200){throw new Error("Bad Server Response")}return result.text()}).then((response)=>{var resp=JSON.parse(response);var nameresp=document.getElementById("name");var messageresp=document.getElementById("message");if(location.protocol == "https:"){url='wss://'+window.location.hostname+'/check?id='+resp.id}else{url='ws://'+window.location.hostname+'/check?id='+resp.id}nameresp.value=resp.name;nameresp.readOnly=true;if(resp.message==""){messageresp.remove();document.getElementById("messagelabel").remove()}else{messageresp.value=resp.message;messageresp.readOnly=true}document.getElementById("amount").readOnly=true;document.getElementById("amount").value=resp.amount;var addrlabel=document.createElement('label');addrlabel.innerText="Address";addrlabel.setAttribute("id","addresslabel");var smallwait=document.createElement('small');smallwait.innerText="Waiting for payment";var addr=document.createElement('blockquote');addr.innerText=resp.address;addr.setAttribute("id","address");var a=document.createElement('a');var linkText=document.createTextNode("Open in GUI");a.appendChild(linkText);a.setAttribute("id","opengui");a.href="monero://"+resp.address+"?tx_amount="+resp.amount;document.getElementById('form').appendChild(addrlabel);document.getElementById('form').appendChild(addr);document.getElementById('form').appendChild(a);c=new WebSocket(url);c.onmessage=function(msg){if(msg.data!="KEEPALIVE"){var wsdata=JSON.parse(msg.data);console.log(wsdata.message);document.getElementById("waiting").remove();document.getElementById("opengui").remove();document.getElementById("address").remove();document.getElementById("addresslabel").remove();document.getElementById("qr").remove();var receipt=document.createElement('blockquote');receipt.innerText=wsdata.message+" XMR received! Superchat sent.";document.getElementById('form').appendChild(receipt)}};var qr=document.createElement('img');qr.src='data:image/png;base64,'+resp.qr;qr.setAttribute("id","qr");document.getElementById('form').appendChild(qr);console.log(JSON.parse(response))}).catch((error)=>{console.log(error)});return false}`

var ConfigJSON string = `{
    "RPCURL": "http://127.0.0.1:18082/json_rpc",
    "DBName": "shadowchat_db",
    "DBHost": "127.0.0.1:5432",
    "DBUser": "shadowchat",
    "DBPass": "",
    "MempoolCacheIntervalMiliseconds": 3000,
    "MoneroGetPriceIntervalSeconds": 480,
    "DeleteUnpaidChatsAfterMinutes": 20,
    "MinimumDonationXMR": 0.0015,
    "MaximumMessageChars": 120,
    "StreamlabsAccessToken": "",
    "StreamlabsCustomImage": "https://www.getmonero.org/press-kit/symbols/monero-symbol-on-white-480.png",
    "StreamlabsCustomSound": "https://powerchat.live/static/audio/coin.mp3",
    "StreamlabsAlertInUSD": false,
	"CookieSigningKey": "YOUMUSTCHANGETHISTOSOMETHINGRANDOM"
}`

var DonoGoalWidgetHTML string = `<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8" />
    <title>Donation Goal</title>
    <link rel="stylesheet" href="/style.css">
    <meta http-equiv="refresh" content="30">
</head>

<style>
body { 
  background: none;
  border: 1px solid #262b3b;
  border-radius: 5px; 
  text-shadow: 1px 2px 3px black; 
  font-weight: bold;
  zoom: 150%;
}
.prog {
  width: calc(({{.Received}}/{{.Goal}})*100%);
  z-index: 8;
}
.r {
  top: 7px; right: 20px;
  z-index: 9;
  position: fixed;
}
.l {
  top: 7px; left: 20px;
  z-index: 9;
  position: fixed;
}
.xmr {
  box-shadow: 1px 1px 2px black; 
  display: inline;
  width: 12px;

}
</style>

<body>
  <blockquote id="prog" class="prog"></blockquote>
  <p class="r">{{.Goal}} XMR</p>
  <p class="l">
    <svg class=xmr version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
    	 viewBox="0 0 512 512" style="enable-background:new 0 0 512 512;" xml:space="preserve">
    <g>
    	<path fill="white" d="M256,0C114.4,0,0,114.6,0,256.3c0,28.5,4.8,55.6,13.2,81.2h76.3V122.2L256,289l166.5-166.8v215.3h76.3
    		c8.2-25.7,13.2-52.7,13.2-81.2C512,114.7,397.6,0,256,0L256,0z M217.8,326.6l-72.8-73v135.5H37.5C82.6,462.7,163.8,512,256,512
    		c92.2,0,174.1-49.3,218.6-123H367V253.6l-72.2,72.9l-38.2,38.3L218,326.5h-0.2L217.8,326.6z"/>
    </g>
    </svg> {{.Received}} 
  </p>
<script>
//function move() {
//  var elem = document.getElementById("prog");   
//  var width = 20;
//  var id = setInterval(frame, 10);
//  function frame() {
//    if (width >= 100) {
//      clearInterval(id);
//    } else {
//      width++; 
//      elem.style.width = width + '%'; 
//      elem.innerHTML = width * 1  + '%';
//    }
//  }
//}
</script>
</body>

</html>`

var IndexHTML string = `<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8" />
    <title>SHDWCHAT</title>
    <link rel="stylesheet" href="/style.css">
    <script src="/scapi.js"></script>
</head>

<body id=body>
    <h1 class=center>SHADOWCHAT</h1>
    <div class=form id=form>
        <form onsubmit="return send()">
            <label for="name">Name:</label>
            <input type="text" id="name" placeholder="Anonymous" />
            <label for="amount" id="amountlabel">Amount:</label>
            <a href="#" id="three">$3</a><a href="#" id="five">$5</a><a href="#" id="ten">$10</a>
            <input type="number" id="amount" step="0.0001" />
            <label for="message" id="messagelabel">Message:</label>
            <textarea id="message" name="message" rows="6"></textarea>
            <input type="submit" id="pay" value="Pay" /></form>
    </div>
    <div id="waiting" class="waiting"><div></div><div>
            <svg version="1.1" x="0px" y="0px" viewBox="0 0 512 512" style="enable-background:new 0 0 512 512;" xml:space="preserve"><g>
            <path fill="#70799b" d="M256,0C114.4,0,0,114.6,0,256.3c0,28.5,4.8,55.6,13.2,81.2h76.3V122.2L256,289l166.5-166.8v215.3h76.3
            c8.2-25.7,13.2-52.7,13.2-81.2C512,114.7,397.6,0,256,0L256,0z M217.8,326.6l-72.8-73v135.5H37.5C82.6,462.7,163.8,512,256,512
            c92.2,0,174.1-49.3,218.6-123H367V253.6l-72.2,72.9l-38.2,38.3L218,326.5h-0.2L217.8,326.6z"/>
            </g></svg></div>
    </div>
    <small>Get wallet: <a href="https://www.getmonero.org/downloads">Desktop</a> - <a href="https://monero.com/wallets">Mobile</a></small>
    <small>Use Cake Wallet for IOS/Android for fast confirmations</small><br>
</body>

</html>`

function getip() {
      var script = document.createElement('script');
      script.type = "text/javascript";
      script.onload = function() {
        document.getElementById("userip").innerHTML=remoteIp;
      };
      script.src = "/remoteIp";
      document.getElementsByTagName('head')[0].appendChild(script);
};

function sendmsg() {
  document.getElementById("results").innerHTML="";
  document.getElementById("results").style.display = "block"
  document.getElementById("submit").textContent="Loading...";
  document.getElementById("submit").disabled = true;

  var socket = new WebSocket('ws://'+window.location.host+"/ws");

  socket.onopen = function() {
    console.log("Socket onopen");
    document.getElementById("results").value='';
    cmd = document.getElementById("cmd").value;
    hostbox= document.getElementById("hostbox").value;
    var msg = {
      cmd: document.getElementById("cmd").value,
      host: document.getElementById("hostbox").value
    };
    socket.send(JSON.stringify(msg));
  };

  socket.onmessage = function (msg) {
    console.log("Socket onmessage");
    if(msg.data!="done")
      document.getElementById("results").append(msg.data+"\n");
    else
      socket.close();
  }

  socket.onclose = function () {
      console.log("Socket onclose");
      document.getElementById("submit").textContent="Run Test";
      document.getElementById("submit").disabled = false;
  }
  return socket;
}

function updateui() {
  document.getElementById("network-title").innerHTML=data.hosts[0].title;
  document.getElementById("network-location").innerHTML=data.hosts[0].location;
  document.getElementById("ipv4").innerHTML=data.hosts[0].ipv4;
  document.getElementById("ipv6").innerHTML=data.hosts[0].ipv6;
  data.hosts[0].testfiles.forEach(function(entry){
      var a = document.createElement('a');
      var linkText = document.createTextNode(entry.size);
      a.appendChild(linkText);
      a.title = entry.size;
      a.href = entry.url;
      document.getElementById("test-files").appendChild(a);
  });
  // closeSpan.setAttribute("class","sr-only");
 // document.getElementById("information").class="span12";

}
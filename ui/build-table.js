function triggerSimpleCapability(device_id, attribute_key, capability_key) {
  $.post("/capability/" + device_id + "/" + attribute_key + "/" + capability_key)
}

function genFullTable(data) {
  let table = document.querySelector("table");
  console.log(data)

  let thead = table.createTHead();
  let row = thead.insertRow(); //Header row
  let th = document.createElement("th");
  th.appendChild(document.createTextNode("Device ID"));
  row.appendChild(th);
  th.appendChild(document.createTextNode("active"));
  row.appendChild(th);

  for (let device of data) {
    let row = table.insertRow();
    row.insertCell().appendChild(document.createTextNode(device["id"]))
    let cell = row.insertCell()
    cell.id = device["id"] + "_active_state"
    cell.appendChild(document.createTextNode(device["attributes"]["active"]["boolean-state"]))

    // Create on button
    let onButton = document.createElement('button');
    onButton.innerHTML = "Activate"
    onButton.addEventListener('click', function(){
      triggerSimpleCapability(device["id"], "active", "activate");
    });
    row.insertCell().appendChild(onButton)

    // Create off button
    let offButton = document.createElement('button');
    offButton.innerHTML = "Deactivate"
    offButton.addEventListener('click', function(){
      triggerSimpleCapability(device["id"], "active", "deactivate");
    });
    row.insertCell().appendChild(offButton)
  }

}
$.getJSON("/discovery", "", genFullTable);
var eSource = new EventSource("/subscribe");

//Now bind various Events , Message, and Error to this event
eSource.addEventListener('open', function(e) {
  console.log("Connection was opened.")
}, false);

eSource.addEventListener('message', function(e) {
  let json_obj = JSON.parse(e.data)
  console.log(json_obj);
  document.getElementById(json_obj.id + "_active_state").innerHTML = json_obj["attributes"]["active"]["boolean-state"]
}, false);

eSource.addEventListener('error', function(e) {
  if (e.readyState == EventSource.CLOSED) {
    console.log("Connection was closed. ");
  }
  console.log(e);
}, false);

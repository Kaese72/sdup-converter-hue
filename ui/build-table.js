function triggerSimpleCapability(device_id, attribute_key, capability_key) {
  $.post("/capability/" + device_id + "/" + attribute_key + "/" + capability_key)
}

function colorchange(device_id, x, y) {
  $.post("/capability/" + device_id + "/color/setkeyval", JSON.stringify({"x": parseFloat(x), "y": parseFloat(y)}), null, "json");
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


    if ('color' in device["attributes"]) {
      // Create color form button
      //let colorform = document.createElement('form');

      let xInput = document.createElement("input")
      xInput.setAttribute("type", "number")
      xInput.setAttribute("value", device["attributes"]["color"]["keyval-state"]["x"])
      xInput.setAttribute("name", "x")
      xInput.id = device["id"] + "_color_x"

      let yInput = document.createElement("input")
      yInput.setAttribute("type", "number")
      yInput.setAttribute("value", device["attributes"]["color"]["keyval-state"]["y"])
      yInput.setAttribute("name", "y")
      yInput.id = device["id"] + "_color_y"

      var submission = document.createElement("button");
      submission.innerHTML = "Change color";
      submission.addEventListener('click', function(){
        colorchange(device["id"], document.getElementById(device["id"] + "_color_x").value, document.getElementById(device["id"] + "_color_y").value);
      });


      //colorform.appendChild(xInput)
      //colorform.appendChild(yInput)
      //colorform.appendChild(submission)
      row.insertCell().appendChild(xInput)
      row.insertCell().appendChild(yInput)
      row.insertCell().appendChild(submission)


      //row.insertCell().appendChild(colorform)
    }
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
  if ('active' in json_obj["attributes"]) {
    document.getElementById(json_obj.id + "_active_state").innerHTML = json_obj["attributes"]["active"]["boolean-state"]

  } else if ('color' in json_obj["attributes"]) {
    document.getElementById(json_obj.id + "_color_y").value = json_obj["attributes"]["color"]["keyval-state"]["y"]
    document.getElementById(json_obj.id + "_color_x").value = json_obj["attributes"]["color"]["keyval-state"]["x"]
  } else {
    console.log("Message does not contain recognizable change")
  }

}, false);

eSource.addEventListener('error', function(e) {
  if (e.readyState == EventSource.CLOSED) {
    console.log("Connection was closed. ");
  }
  console.log(e);
}, false);

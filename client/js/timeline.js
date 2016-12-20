var timeline;

function orderFunc(a, b) {
  return a.subgroupOrder - b.subgroupOrder;
}

function genGroups(numGroups) {
  var groups = [];
  for (let i = 0; i < numGroups; i++) {
    var group = {
      id: i,
      content: "CU " + i,
      value: i,
      subgroupOrder: orderFunc
    };
    groups.push(group);
  }
  return groups;
}

function timelineFunc() {
  var filterForm = document.getElementById('filterForm').value;
  var query = 'timeline/json?' + filterForm;
  console.log(query);
  axios.get(query)
      .then(function(response) {

        // DOM element where the Timeline will be attached
        var container = document.getElementById('timeline_container');

        // Clear
        while (container.hasChildNodes()) {
           container.removeChild(container.lastChild);
        }

        console.log(response.data)

        // Create a DataSet (allows two way data-binding)
        var items = new vis.DataSet(response.data);

        // Configuration for the Timeline
        var options = {
          stack: false,
          editable: false,
          groupOrder: orderFunc,
          format: {
            minorLabels: {millisecond: "x"},
            majorLabels: {millisecond: "x"}
          },
          min: 0
        };

        // Create a Timeline
        timeline = new vis.Timeline(container, items, options);
        var groups = genGroups(1);
        timeline.setGroups(new vis.DataSet(groups));

      })
      .catch(function(error) { console.log(error); });
}

function moveTimeline(percentage) {
  var range = timeline.getWindow();
  var interval = range.end - range.start;
  timeline.setWindow({
    start: range.start.valueOf() - interval * percentage,
    end: range.end.valueOf() - interval * percentage
  });
}

document.addEventListener("keydown", function(event) {
  var percentage = 0.2;
  if (event.ctrlKey) {
    percentage = 0.8
  }

  switch (event.keyCode) {
    case 37:
      moveTimeline(percentage);
      break;
    case 39:
      moveTimeline(-percentage);
      break;
    default:
      break;
  }
});
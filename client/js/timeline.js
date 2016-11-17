var groups = new vis.DataSet([
    {id: 0, content: "CU 0", value: 0, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 1, content: "CU 1", value: 1, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 2, content: "CU 2", value: 2, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 3, content: "CU 3", value: 3, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 4, content: "CU 4", value: 4, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 5, content: "CU 5", value: 5, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 6, content: "CU 6", value: 6, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 7, content: "CU 7", value: 7, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 8, content: "CU 8", value: 8, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 9, content: "CU 9", value: 9, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 10, content: "CU 10", value: 10, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 11, content: "CU 11", value: 11, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 12, content: "CU 12", value: 12, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 13, content: "CU 13", value: 13, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 14, content: "CU 14", value: 14, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 15, content: "CU 15", value: 15, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 16, content: "CU 16", value: 16, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 17, content: "CU 17", value: 17, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 18, content: "CU 18", value: 18, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 19, content: "CU 19", value: 19, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 20, content: "CU 20", value: 20, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 21, content: "CU 21", value: 21, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 22, content: "CU 22", value: 22, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 23, content: "CU 23", value: 23, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 24, content: "CU 24", value: 24, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 25, content: "CU 25", value: 25, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 26, content: "CU 26", value: 26, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 27, content: "CU 27", value: 27, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 28, content: "CU 28", value: 28, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 29, content: "CU 29", value: 29, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 30, content: "CU 30", value: 30, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
    {id: 31, content: "CU 31", value: 31, subgroupOrder: function (a, b) {return a.subgroupOrder - b.subgroupOrder;}},
]);


axios.get('/timeline/json')
  .then(function (response) {

    // DOM element where the Timeline will be attached
    var container = document.getElementById('timeline_container');

    // Create a DataSet (allows two way data-binding)
    var items = new vis.DataSet(response.data);

    // Configuration for the Timeline
    var options = {
        stack: false,
        editable: false,
        groupOrder: function (a, b) {
          return a.value - b.value;
        },
        format: {
            minorLabels: {millisecond: "x"},
            majorLabels: {millisecond: "x"}
        },
        min: 0
    };

    // Create a Timeline
    var timeline = new vis.Timeline(container, items, options);
    timeline.setGroups(groups);

  })
  .catch(function (error) {
    console.log(error);
  });


$(function() {
	var ctx = $("#ants")[0].getContext("2d");


	$.get( "http://localhost:8081/api/world", function( data ) {
		var len = data.WorldMap.length;
		ctx.fillStyle = "#000000";
		ctx.fillRect(0,0,len,len);
		for (var i = 0; i < len; i++) {
			var row_len = data.WorldMap[i].length;
			for (var j = 0; j < row_len; j++) {
				if (data.WorldMap[i][j].Signatures.search < 10000000) {
					ctx.fillStyle = "#550000";
					ctx.fillRect(i,j,1,1);
				}
				if (data.WorldMap[i][j].Signatures.food < 10000000) {
					ctx.fillStyle = "#FFFFFF";
					ctx.fillRect(i,j,1,1);
				}
				if (data.WorldMap[i][j].Food > 0) {
					ctx.fillStyle = "#00FF00";
					ctx.fillRect(i,j,10,10);
				}
			}
		}
	}, "json");


});

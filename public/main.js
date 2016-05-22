$(function() {

    /**
     * WebRTC.
     */
    var streaming = false;
    var width = 320;
    var height = 0;
    var dtpic = 2000;

    var video;
    var canvas;
    var photo;
    var startButton;

    var startVideo = function() {
      video = document.getElementById('video');
      canvas = document.getElementById('canvas');
      photo = document.getElementById('photo');
      startbutton = document.getElementById('startbutton');

      navigator.getMedia = ( navigator.getUserMedia ||
                             navigator.webkitGetUserMedia ||
                             navigator.mozGetUserMedia ||
                             navigator.msGetUserMedia);

      if (!navigator.getMedia) {
        alert('Browser does not support webRTC.');
      }

      navigator.getMedia({
          video: true,
          audio: false
        },
        function(stream) {
          console.log('Streaming...');
          if (navigator.mozGetUserMedia) {
            video.mozSrcObject = stream;
          } else {
            var vendorURL = window.URL || window.webkitURL;
            video.src = vendorURL.createObjectURL(stream);
          }
          video.play();
        },
        function(err) {
            console.log("An error occured! ", err);
        }
      );

      video.addEventListener('canplay', function(ev){
        if (!streaming) {
          height = video.videoHeight / (video.videoWidth/width);

          if (isNaN(height)) {
            height = width / (4/3);
          }

          video.setAttribute('width', width);
          video.setAttribute('height', height);
          canvas.setAttribute('width', width);
          canvas.setAttribute('height', height);
          streaming = true;
        }
      }, false);

    }

    window.addEventListener('load', startVideo, false);

    /**
     * Websockets.
     */
    if (!window["WebSocket"]) {
        alert('Browser does not support WebSocket.');
    }

    var content = $("#content");
    var conn = new WebSocket('ws://' + window.location.host + '/ws');

    conn.onopen = function() {
        console.log('Connection open.');
    };

    conn.onclose = function() {
        console.log('Connection closed.');
    };

    conn.onmessage = function(msg) {
        var message = JSON.parse(msg.data);
        var tags = message.Tags;
        var pic = message.Pic;
        if (tags && tags.length > 0) {
          var photo = document.getElementById('photo');
          photo.src = pic;
          var headerMessage = 'Last tag' + ((tags.length > 1) ? 's' : '') + ' detected: ' + tags.join(', ');
          $('#tags-header').html(headerMessage);
        }
    };

    /**
     * Webcam shot.
     */
    var getPicture = function() {
      var context = canvas.getContext('2d');
      if (width && height) {
        canvas.width = width;
        canvas.height = height;
        context.drawImage(video, 0, 0, width, height);

        var data = canvas.toDataURL('image/png');
        return data;
      }
    }

    var sendPictureAndTags = function(tags) {
      var pic = getPicture();
      var message = {
        type: 'picture',
        pic: pic,
        tags: tags
      };
      conn.send(JSON.stringify(message));
    }

    var picInterval;
    $('#stop').prop('disabled', true);

    /**
     * Form.
     */
    $('#submit').click(function(ev) {
      ev.preventDefault();
      var tags = [];
      $('input[name=tags]').val().trim().split(',').forEach(function(str) {
        var s = str.trim();
        if (s && s !== '') {
          tags.push(s);
        }
      });
      if (picInterval) {
        clearInterval(picInterval);
      }

      if (tags && tags.length > 0) {
        $('#stop').prop('disabled', false);
        picInterval = setInterval(function() {
          sendPictureAndTags(tags);
        }, dtpic);
      } else {
        alert('Enter some tags to be detected.');
      }
    });

    $('#stop').click(function(ev) {
      ev.preventDefault();
      $('#stop').prop('disabled', true);
      if (picInterval) {
        clearInterval(picInterval);
      }
    });


});

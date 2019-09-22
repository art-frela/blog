$(document).ready(function () {

});

let apiPostURL = "/api/v1/posts"
let userID = "00000000-0000-0000-00000000"

// events listeners
$('.saveeditpost').bind('click', function(e){
    var id = $(this).attr("task-id")
    var title = $(".post_title_edit").val()
    var rubric_id = $('.post_rubric_edit :selected').val()
    var content = $(".post_content_edit").val()
    newUpdPost(title, rubric_id, content, 'put', id)
    e.stopPropagation()
})

$('.savenewpost').bind('click', function(e){
    var id = $(this).attr("task-id")
    var title = $(".post_title_edit").val()
    var rubric_id = $('.post_rubric_edit :selected').val()
    var content = $(".post_content_edit").val()
    newUpdPost(title, rubric_id, content, 'post', id)
    e.stopPropagation()
})


// functions
function newUpdPost(title, rubric_id, content, method, id) {
    var data = {
        title: title,
        content: content,
        user_id: userID,
        rubric_id: rubric_id
    };
    var url = apiPostURL
    if (method == 'put') {
        url += "/"+id
    }
    $.ajax({
        url: url,
        cache: false,
        type: method,
        data: JSON.stringify(data),
        headers: {
            "Content-type": "application/json"
        },
        success: function (html) {

            console.dir(html)
            var locURL = "/posts/"
            if (method == 'put') {
                locURL += id
            } else {
                locURL += html.message
            }
            document.location = locURL
            //result = $.parseJSON(html);
        },
        error: function (request, status, error) {
            console.error(request+"; "+status+"; "+error)
        }
    });
}
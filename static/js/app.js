document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault();  // 阻止表单默认提交行为

    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;

    // 构建发送到服务器的数据对象
    var data = {
        username: username,
        password: password
    };

    // 使用 fetch API 发送数据到后端服务器
    fetch('/register', {  // 假设你的服务器有一个名为 '/register' 的端点来处理注册请求
        method: 'POST',   // 设置请求方法为 POST
        headers: {
            'Content-Type': 'application/json'  // 设置请求头，告诉服务器消息主体是 JSON 格式
        },
        body: JSON.stringify(data)  // 将 JavaScript 对象转换成 JSON 字符串
    })
        .then(function(response) {
            return response.json();  // 解析服务器返回的 JSON 数据
        })
        .then(function(data) {
            console.log('Success:', data);
            window.location.href = "/"
            alert("注册成功！用户名: " + username);
        })
        .catch(function(error) {
            console.error('Error:', error);
            alert("注册失败，请重试！");
        });
});

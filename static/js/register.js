function is_same() {
	psw = document.getElementById("register-psw")
	cpsw = document.getElementById("register-cpsw")
	if (psw.value === cpsw.value) {
		return true
	}
	return false
}

function check_form() {
	help = document.getElementById("cpsw-help")
	btn = document.getElementById("register-btn")
	//如果密码相同的话
	if (is_same()) {
		help.innerHTML = "😄密码一致，可以注册。"
		btn.disabled = false
	} else {
		help.innerHTML = "🙃密码不一致！"
		btn.disabled = true
	}
}
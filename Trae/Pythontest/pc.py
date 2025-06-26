from seleniumbase import SB
import json
import os

def save_cookies(driver, filename='cookies.json'):
    cookies = driver.get_cookies()
    with open(filename, 'w') as f:
        json.dump(cookies, f)

def load_cookies(driver, filename='cookies.json'):
    if os.path.exists(filename):
        with open(filename, 'r') as f:
            cookies = json.load(f)
            for cookie in cookies:
                driver.add_cookie(cookie)
        return True
    return False

# 基本用法
with SB() as sb:
    # 打开登录页面
    sb.open("https://link.nertron.net/#/login")
    
    # 尝试加载已保存的cookies
    if load_cookies(sb.driver):
        # 刷新页面以应用cookies
        sb.refresh()
        sb.sleep(2)
        
        # 检查是否已登录（检查URL或特定元素）
        if not "login" in sb.get_current_url():
            print("已通过cookies自动登录")
        else:
            print("cookies已失效，需要重新登录")
    if "login" in sb.get_current_url():
        # 等待页面加载完成
        sb.wait_for_element('input[placeholder="请输入邮箱"]')
        # 输入邮箱
        sb.type('input[placeholder="请输入邮箱"]', 'zhiao321@gmail.com')
        # 输入密码
        sb.type('input[placeholder="请输入密码"]', 'Theo@1223')
        # 点击登录按钮
        sb.click('button.ant-btn-primary')
        # 等待登录完成
        sb.sleep(3)
        
        # 检查是否存在购买订阅按钮来验证登录状态
        if sb.is_element_visible('span.tit:contains("购买订阅")'):
            save_cookies(sb.driver)
            print("登录成功，已保存cookies")
        else:
            print("登录可能失败，请检查账号密码")
    
    # 等待并点击购买订阅按钮
    sb.wait_for_element('span.tit:contains("购买订阅")')
    sb.click('span.tit:contains("购买订阅")')
    
    
    # 等待操作完成
    sb.sleep(2)
    
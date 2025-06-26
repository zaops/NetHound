from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
from bs4 import BeautifulSoup
import time

# 设置 Chrome 浏览器驱动
driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()))

# 访问目标网页
url = "https://vacations.ctrip.com/travel/detail/p35548034?city=1"  # 替换为实际 URL
driver.get(url)

# 等待页面加载（根据需要调整等待时间）
time.sleep(3)

# 获取渲染后的页面源码
soup = BeautifulSoup(driver.page_source, "html.parser")

# 查找用户名
user_elements = soup.find_all("p", class_="ct-review-person-name")
for user in user_elements:
    username = user.get_text(strip=True)
    print(f"用户名: {username}")

# 关闭浏览器
driver.quit()


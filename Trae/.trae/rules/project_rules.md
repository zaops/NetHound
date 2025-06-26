# SeleniumBase 完整使用教程

## 目录

1. [简介](#简介)
2. [基本概念](#基本概念)
3. [基础用法](#基础用法)
   - [初始化与浏览器控制](#初始化与浏览器控制)
   - [页面交互](#页面交互)
   - [元素定位](#元素定位)
   - [等待机制](#等待机制)
   - [断言验证](#断言验证)
4. [高级功能](#高级功能)
   - [CDP 模式（反爬虫）](#cdp-模式反爬虫)
   - [UC 模式（无头浏览器）](#uc-模式无头浏览器)
   - [JavaScript 执行](#javascript-执行)
   - [文件上传与下载](#文件上传与下载)
   - [处理弹窗和框架](#处理弹窗和框架)
   - [多浏览器并行](#多浏览器并行)
5. [爬虫应用场景](#爬虫应用场景)
   - [基本数据抓取](#基本数据抓取)
   - [绕过验证码和人机验证](#绕过验证码和人机验证)
   - [处理动态加载内容](#处理动态加载内容)
   - [模拟用户行为](#模拟用户行为)
   - [数据提取与保存](#数据提取与保存)
6. [测试框架集成](#测试框架集成)
   - [与 pytest 集成](#与-pytest-集成)
   - [测试报告生成](#测试报告生成)
   - [并行测试执行](#并行测试执行)
7. [常见问题与解决方案](#常见问题与解决方案)
8. [最佳实践](#最佳实践)

## 简介

SeleniumBase 是一个功能强大的 Python 自动化框架，它在 Selenium WebDriver 的基础上进行了封装和扩展，提供了更简洁、更可靠的 API，使网页自动化操作变得更加简单高效。SeleniumBase 不仅适用于网页测试，还非常适合网页爬虫、数据抓取和自动化任务。

与原生 Selenium 相比，SeleniumBase 具有以下优势：

1. **简化的 API**：SeleniumBase 的方法通常会执行多个操作。例如，`type()` 方法会自动等待元素可见、清除文本框、输入文本，如果文本以 `\n` 结尾还会自动提交。

2. **智能等待**：SeleniumBase 内置了智能等待机制，不需要手动添加显式等待，大大减少了代码量和复杂度。

3. **自动驱动管理**：自动下载和管理与浏览器匹配的驱动程序，无需手动配置。

4. **选择器自动检测**：自动识别 CSS 选择器和 XPath，无需指定选择器类型。

5. **反爬虫能力**：提供 CDP 模式和 UC 模式，能够绕过许多网站的反爬虫机制和验证码。

6. **丰富的调试功能**：提供截图、日志记录、演示模式等调试功能。

7. **测试框架集成**：与 pytest、nose、behave 等测试框架无缝集成。

## 基本概念

在开始使用 SeleniumBase 之前，了解一些基本概念是很有帮助的：

**SB 类**：SeleniumBase 提供了一个名为 `SB` 的类，它是一个上下文管理器，可以使用 `with` 语句简化浏览器的创建和关闭。

**BaseCase 类**：这是 SeleniumBase 的核心类，包含了大量用于浏览器自动化的方法。当与测试框架如 pytest 一起使用时，你的测试类需要继承这个类。

**选择器**：SeleniumBase 支持多种元素定位方式，包括 CSS 选择器、XPath、ID、名称、类名等。大多数情况下，CSS 选择器是最推荐的方式。

**超时设置**：SeleniumBase 的方法都有默认的超时设置，可以根据需要进行调整。

## 基础用法

### 初始化与浏览器控制

SeleniumBase 提供了两种主要的使用方式：使用 `SB` 类或继承 `BaseCase` 类。

#### 使用 SB 类（推荐用于爬虫）

```python
from seleniumbase import SB

# 基本用法
with SB() as sb:
    sb.open("https://www.example.com")
    # 执行其他操作...

# 使用更多选项
with SB(uc=True, headless=False, test=True) as sb:
    sb.open("https://www.example.com")
    # 执行其他操作...
```

`SB` 类的常用参数：

- `browser`：浏览器类型，默认为 "chrome"，也支持 "firefox"、"edge" 等
- `headless`：是否使用无头模式，默认为 False
- `uc`：是否使用 UC 模式（反检测模式），默认为 False
- `test`：是否在测试模式下运行，默认为 False
- `incognito`：是否使用隐身模式，默认为 False
- `locale`：浏览器语言设置，如 "zh-CN"、"en-US" 等
- `proxy`：代理服务器设置，如 "socks5://127.0.0.1:1080"

#### 继承 BaseCase 类（推荐用于测试）

```python
from seleniumbase import BaseCase
BaseCase.main(__name__, __file__)  # 调用 pytest

class MyTestClass(BaseCase):
    def test_example(self):
        self.open("https://www.example.com")
        # 执行其他操作...
```

### 页面交互

SeleniumBase 提供了丰富的方法来与网页元素进行交互：

#### 打开网页

```python
# 使用 SB 类
with SB() as sb:
    sb.open("https://www.example.com")

# 使用 BaseCase 类
class MyTest(BaseCase):
    def test_example(self):
        self.open("https://www.example.com")
```

#### 点击元素

```python
# 点击按钮
sb.click("button#submit")  # 使用 CSS 选择器
sb.click("//button[@id='submit']")  # 使用 XPath

# 点击链接
sb.click("a[href='/login']")

# 设置超时时间
sb.click("button#submit", timeout=10)  # 等待最多 10 秒
```

#### 输入文本

```python
# 在输入框中输入文本
sb.type("#username", "user123")
sb.type("#password", "password123")

# 输入后自动提交（按回车键）
sb.type("#search", "SeleniumBase tutorial\n")

# 清除输入框后输入
sb.clear("#username")
sb.type("#username", "new_user")

# 或者直接使用 type 方法（会自动清除）
sb.type("#username", "new_user")
```

#### 选择下拉菜单选项

```python
# 通过可见文本选择
sb.select_option_by_text("select#country", "China")

# 通过值选择
sb.select_option_by_value("select#country", "CN")

# 通过索引选择
sb.select_option_by_index("select#country", 1)  # 选择第二个选项（索引从 0 开始）
```

#### 拖放操作

```python
# 拖放元素
sb.drag_and_drop("div#source", "div#target")
```

#### 滚动页面

```python
# 滚动到元素可见
sb.scroll_to("footer")

# 滚动到页面顶部
sb.scroll_to_top()

# 滚动到页面底部
sb.scroll_to_bottom()
```

#### 悬停操作

```python
# 鼠标悬停在元素上
sb.hover_and_click("a#menu", "a#submenu")  # 悬停在菜单上，然后点击子菜单
sb.hover("div.tooltip")  # 仅悬停
```

### 元素定位

SeleniumBase 支持多种元素定位方式，并且会自动检测选择器类型：

#### CSS 选择器（推荐）

```python
sb.click("button#submit")  # ID 选择器
sb.click("button.primary")  # 类选择器
sb.click("div.header button")  # 后代选择器
sb.click("input[name='username']")  # 属性选择器
sb.click("ul > li:first-child")  # 子元素和伪类选择器
```

#### XPath 选择器

```python
sb.click("//button[@id='submit']")
sb.click("//div[@class='header']//button")
sb.click("//input[@name='username']")
sb.click("//ul/li[1]")
```

#### 其他定位方式

```python
# 通过链接文本
sb.click_link("Sign In")

# 通过部分链接文本
sb.click_link_text("Sign")

# 通过 JavaScript
sb.js_click("#submit-button")  # 当常规点击不起作用时使用
```

#### 查找多个元素

```python
# 获取所有匹配的元素
elements = sb.find_elements("div.product")

# 遍历元素
for element in elements:
    name = sb.get_text(element.find_element_by_css_selector("h3"))
    price = sb.get_text(element.find_element_by_css_selector(".price"))
    print(f"Product: {name}, Price: {price}")
```

### 等待机制

SeleniumBase 的大多数方法都内置了智能等待机制，但有时你可能需要显式等待特定条件：

#### 基本等待

```python
# 等待指定秒数
sb.sleep(2)  # 等待 2 秒

# 等待元素出现
sb.wait_for_element("div#loading", timeout=10)  # 最多等待 10 秒

# 等待元素消失
sb.wait_for_element_absent("div#loading", timeout=10)

# 等待文本出现
sb.wait_for_text("Welcome", "div#message", timeout=5)
```

#### 高级等待

```python
# 等待元素可见
sb.wait_for_element_visible("button#submit")

# 等待元素可点击
sb.wait_for_element_clickable("button#submit")

# 等待页面加载完成
sb.wait_for_ready_state_complete()

# 等待 URL 包含特定文本
sb.wait_for_url_to_contain("dashboard")
```

### 断言验证

SeleniumBase 提供了丰富的断言方法，用于验证页面状态：

```python
# 验证元素存在
sb.assert_element("div.success")

# 验证元素不存在
sb.assert_element_absent("div.error")

# 验证文本存在
sb.assert_text("Welcome", "h1")

# 验证文本不存在
sb.assert_text_not_visible("Error", "div#message")

# 验证属性值
sb.assert_attribute("img#logo", "alt", "Company Logo")

# 验证标题
sb.assert_title("Dashboard - Example Site")

# 验证 URL
sb.assert_url_contains("dashboard")

# 验证元素可见
sb.assert_element_visible("button#submit")

# 验证元素不可见
sb.assert_element_not_visible("div#loading")
```

## 高级功能

### CDP 模式（反爬虫）

CDP（Chrome DevTools Protocol）模式是 SeleniumBase 的一个强大功能，它允许你绕过许多网站的反爬虫机制，包括 Cloudflare 等保护措施。

```python
from seleniumbase import SB

# 基本 CDP 模式用法
with SB(uc=True) as sb:
    url = "https://example.com"
    sb.activate_cdp_mode(url)
    # 如果遇到验证码或挑战页面
    sb.uc_gui_click_captcha()
    # 继续其他操作...
```

更复杂的 CDP 模式示例，用于绕过 Cloudflare 保护：

```python
from seleniumbase import SB

with SB(uc=True, test=True, locale="en") as sb:
    url = "https://site-with-cloudflare.com"
    sb.activate_cdp_mode(url)
    sb.uc_gui_click_captcha()
    sb.sleep(2)  # 等待验证完成
    
    # 验证是否成功绕过保护
    if sb.is_text_visible("Welcome", "h1"):
        print("成功绕过 Cloudflare 保护！")
        
    # 继续抓取数据
    data = sb.get_text("div.content")
    print(data)
```

### UC 模式（无头浏览器）

UC（Undetected Chrome）模式是另一种反检测技术，它修改了 Chrome 的特征，使其更难被网站识别为自动化浏览器：

```python
from seleniumbase import SB

# 使用 UC 模式
with SB(uc=True) as sb:
    sb.open("https://bot.sannysoft.com")  # 一个检测浏览器自动化的网站
    # 查看是否被检测为机器人
    detection_result = sb.get_text("body")
    print("检测结果:", detection_result)
    sb.save_screenshot("bot_detection_result.png")
```

### JavaScript 执行

SeleniumBase 允许你执行 JavaScript 代码，这对于某些复杂操作或绕过前端限制非常有用：

```python
# 执行简单的 JavaScript
sb.execute_script("return document.title;")

# 修改页面元素
sb.execute_script("document.getElementById('header').style.display = 'none';")

# 滚动到元素
sb.execute_script("arguments[0].scrollIntoView();", sb.find_element("div#target"))

# 获取页面数据
data = sb.execute_script("""
    let results = [];
    document.querySelectorAll('.product').forEach(product => {
        results.push({
            name: product.querySelector('h3').innerText,
            price: product.querySelector('.price').innerText,
            url: product.querySelector('a').href
        });
    });
    return results;
""")
print(data)
```

### 文件上传与下载

处理文件上传：

```python
# 上传文件
sb.choose_file("input#fileUpload", "/path/to/file.jpg")

# 验证上传成功
sb.assert_text("File uploaded successfully", "div.message")
```

处理文件下载：

```python
# 设置下载目录
download_dir = "/path/to/downloads"
with SB(browser="chrome", downloads_folder=download_dir) as sb:
    sb.open("https://example.com/download-page")
    sb.click("a#downloadButton")
    
    # 等待下载完成
    sb.wait_for_file_download("example.pdf", timeout=30)
    
    # 验证文件存在
    import os
    assert os.path.exists(os.path.join(download_dir, "example.pdf"))
```

### 处理弹窗和框架

处理 JavaScript 弹窗：

```python
# 处理 alert 弹窗
sb.click("button#showAlert")
sb.accept_alert()  # 点击 "确定"

# 处理 confirm 弹窗
sb.click("button#showConfirm")
sb.dismiss_alert()  # 点击 "取消"
sb.accept_alert()  # 点击 "确定"

# 处理 prompt 弹窗
sb.click("button#showPrompt")
sb.send_keys_to_alert("My input")  # 输入文本
sb.accept_alert()  # 点击 "确定"
```

处理 iframe：

```python
# 切换到 iframe
sb.switch_to_frame("iframe#content")

# 在 iframe 中操作
sb.type("#username", "user123")
sb.type("#password", "pass123")
sb.click("button#login")

# 切回主文档
sb.switch_to_default_content()
```

### 多浏览器并行

在爬虫场景中，有时需要同时操作多个浏览器实例以提高效率：

```python
from seleniumbase import SB
import threading

def scrape_website(url, results, index):
    with SB(uc=True) as sb:
        sb.open(url)
        # 执行抓取操作
        title = sb.get_title()
        content = sb.get_text("div.content")
        results[index] = {"title": title, "content": content}

# 要抓取的网址列表
urls = [
    "https://example.com/page1",
    "https://example.com/page2",
    "https://example.com/page3",
    "https://example.com/page4"
]

# 存储结果
results = [None] * len(urls)

# 创建并启动线程
threads = []
for i, url in enumerate(urls):
    thread = threading.Thread(target=scrape_website, args=(url, results, i))
    threads.append(thread)
    thread.start()

# 等待所有线程完成
for thread in threads:
    thread.join()

# 处理结果
for result in results:
    print(result)
```

## 爬虫应用场景

### 基本数据抓取

以下是一个基本的网页数据抓取示例：

```python
from seleniumbase import SB
import csv

def scrape_products():
    with SB(uc=True) as sb:
        sb.open("https://example.com/products")
        
        # 等待产品列表加载
        sb.wait_for_element("div.product-list")
        
        # 获取所有产品元素
        products = sb.find_elements("div.product")
        
        results = []
        for product in products:
            # 提取产品信息
            name = sb.get_text("h3.title", element=product)
            price = sb.get_text("span.price", element=product)
            description = sb.get_text("p.description", element=product)
            
            results.append({
                "name": name,
                "price": price,
                "description": description
            })
        
        return results

# 执行抓取
products = scrape_products()

# 保存到 CSV 文件
with open("products.csv", "w", newline="", encoding="utf-8") as f:
    writer = csv.DictWriter(f, fieldnames=["name", "price", "description"])
    writer.writeheader()
    writer.writerows(products)

print(f"已抓取 {len(products)} 个产品并保存到 products.csv")
```

### 绕过验证码和人机验证

使用 CDP 模式绕过 Cloudflare 和其他保护：

```python
from seleniumbase import SB
import time

def scrape_protected_site():
    with SB(uc=True, test=True) as sb:
        url = "https://protected-site.com"
        
        # 激活 CDP 模式并访问网站
        sb.activate_cdp_mode(url)
        
        # 处理可能出现的验证码
        if sb.is_element_visible("iframe[title*='challenge']"):
            print("检测到验证码，尝试绕过...")
            sb.uc_gui_click_captcha()
            sb.sleep(3)  # 等待验证完成
        
        # 验证是否成功绕过
        if sb.is_text_visible("Welcome", "h1"):
            print("成功绕过保护！")
            
            # 继续抓取数据
            data = sb.get_text("div.content")
            return data
        else:
            print("绕过保护失败")
            return None

# 执行抓取
content = scrape_protected_site()
if content:
    with open("protected_content.txt", "w", encoding="utf-8") as f:
        f.write(content)
    print("内容已保存到 protected_content.txt")
```

### 处理动态加载内容

许多现代网站使用 JavaScript 动态加载内容，这需要特殊处理：

```python
from seleniumbase import SB
import time
import json

def scrape_infinite_scroll():
    with SB(uc=True) as sb:
        sb.open("https://example.com/infinite-scroll-page")
        
        # 设置要抓取的项目数量
        target_items = 100
        items = []
        
        while len(items) < target_items:
            # 获取当前可见的项目
            current_items = sb.find_elements("div.item")
            
            # 提取数据
            for item in current_items:
                if item not in items:  # 避免重复
                    item_data = {
                        "title": sb.get_text("h3", element=item),
                        "description": sb.get_text("p", element=item)
                    }
                    items.append(item_data)
                    print(f"已抓取 {len(items)} 个项目")
                    
                    if len(items) >= target_items:
                        break
            
            # 滚动到底部加载更多
            sb.execute_script("window.scrollTo(0, document.body.scrollHeight);")
            sb.sleep(2)  # 等待新内容加载
            
            # 检查是否已经到达页面底部且没有更多内容
            if sb.execute_script("return (window.innerHeight + window.scrollY) >= document.body.scrollHeight;"):
                # 尝试点击"加载更多"按钮（如果存在）
                if sb.is_element_visible("button#load-more"):
                    sb.click("button#load-more")
                    sb.sleep(2)
                else:
                    print("已到达页面底部，没有更多内容")
                    break
        
        return items

# 执行抓取
items = scrape_infinite_scroll()

# 保存为 JSON
with open("infinite_scroll_items.json", "w", encoding="utf-8") as f:
    json.dump(items, f, ensure_ascii=False, indent=2)

print(f"已抓取 {len(items)} 个项目并保存到 infinite_scroll_items.json")
```

### 模拟用户行为

有时需要模拟真实用户行为以避免被检测为爬虫：

```python
from seleniumbase import SB
import random
import time

def browse_like_human(url):
    with SB(uc=True) as sb:
        # 打开网页
        sb.open(url)
        
        # 随机等待，模拟阅读
        time.sleep(random.uniform(2, 5))
        
        # 随机滚动
        for _ in range(random.randint(3, 8)):
            # 随机滚动距离
            scroll_amount = random.randint(100, 500)
            sb.execute_script(f"window.scrollBy(0, {scroll_amount});")
            time.sleep(random.uniform(1, 3))
        
        # 随机点击一个链接
        links = sb.find_elements("a[href]")
        if links:
            random_link = random.choice(links)
            sb.hover_and_click(random_link)
            
            # 在新页面上等待和滚动
            time.sleep(random.uniform(3, 7))
            for _ in range(random.randint(2, 6)):
                scroll_amount = random.randint(100, 400)
                sb.execute_script(f"window.scrollBy(0, {scroll_amount});")
                time.sleep(random.uniform(1, 2))
        
        # 获取当前页面内容
        title = sb.get_title()
        content = sb.get_text("body")
        
        return {
            "title": title,
            "url": sb.get_current_url(),
            "content_length": len(content)
        }

# 执行模拟浏览
result = browse_like_human("https://example.com")
print(result)
```

### 数据提取与保存

从复杂网页提取结构化数据并保存：

```python
from seleniumbase import SB
import pandas as pd
import os

def extract_table_data(url, table_selector):
    with SB(uc=True) as sb:
        sb.open(url)
        
        # 等待表格加载
        sb.wait_for_element(table_selector)
        
        # 提取表头
        headers = []
        header_cells = sb.find_elements(f"{table_selector} th")
        for cell in header_cells:
            headers.append(sb.get_text(cell))
        
        # 提取行数据
        rows = []
        table_rows = sb.find_elements(f"{table_selector} tbody tr")
        
        for row in table_rows:
            row_data = []
            cells = sb.find_elements("td", element=row)
            for cell in cells:
                row_data.append(sb.get_text(cell))
            
            if len(row_data) == len(headers):  # 确保数据对齐
                rows.append(row_data)
        
        # 创建 DataFrame
        df = pd.DataFrame(rows, columns=headers)
        return df

# 执行数据提取
url = "https://example.com/data-table"
table_selector = "table.data-table"
df = extract_table_data(url, table_selector)

# 保存为多种格式
os.makedirs("extracted_data", exist_ok=True)

# CSV 格式
df.to_csv("extracted_data/table_data.csv", index=False, encoding="utf-8")

# Excel 格式
df.to_excel("extracted_data/table_data.xlsx", index=False)

# JSON 格式
df.to_json("extracted_data/table_data.json", orient="records", force_ascii=False, indent=2)

print(f"已提取 {len(df)} 行数据并保存到 extracted_data 目录")
```

## 测试框架集成

### 与 pytest 集成

SeleniumBase 与 pytest 无缝集成，这对于创建可维护的自动化测试非常有用：

```python
from seleniumbase import BaseCase
BaseCase.main(__name__, __file__)  # 调用 pytest

class TestExample(BaseCase):
    def test_search_functionality(self):
        # 打开网站
        self.open("https://www.google.com")
        
        # 搜索关键词
        self.type('input[name="q"]', "SeleniumBase tutorial\n")
        
        # 验证搜索结果
        self.assert_element("div#search")
        self.assert_text("SeleniumBase", "div#search")
        
        # 截图
        self.save_screenshot("search_results.png")
```

运行测试：

```bash
pytest test_example.py
```

使用命令行选项：

```bash
# 指定浏览器
pytest test_example.py --browser=firefox

# 使用无头模式
pytest test_example.py --headless

# 使用演示模式（慢速执行并高亮元素）
pytest test_example.py --demo

# 并行执行测试
pytest test_example.py -n=4

# 生成 HTML 报告
pytest test_example.py --html=report.html
```

### 测试报告生成

SeleniumBase 可以生成详细的测试报告：

```python
from seleniumbase import BaseCase
BaseCase.main(__name__, __file__)

class TestReportExample(BaseCase):
    def test_report_example(self):
        self.open("https://www.example.com")
        self.assert_title("Example Domain")
        self.assert_element("h1")
        self.assert_text("Example Domain", "h1")
```

生成 HTML 报告：

```bash
pytest test_report_example.py --html=report.html
```

生成 Dashboard：

```bash
pytest test_report_example.py --dashboard
```

查看 Dashboard：

```bash
pytest test_report_example.py --dashboard --dashboard-show
```

### 并行测试执行

使用 pytest-xdist 插件实现并行测试：

```python
from seleniumbase import BaseCase
BaseCase.main(__name__, __file__)

class TestParallel(BaseCase):
    def test_page_1(self):
        self.open("https://example.com/page1")
        # 测试逻辑...
    
    def test_page_2(self):
        self.open("https://example.com/page2")
        # 测试逻辑...
    
    def test_page_3(self):
        self.open("https://example.com/page3")
        # 测试逻辑...
```

并行执行测试：

```bash
pytest test_parallel.py -n=3  # 使用 3 个并行进程
```

## 常见问题与解决方案

### 元素无法点击

有时元素虽然存在但无法点击，可以尝试以下解决方案：

```python
# 使用 JavaScript 点击
sb.js_click("#element")

# 先滚动到元素，再点击
sb.scroll_to("#element")
sb.click("#element")

# 使用 ActionChains 点击
from selenium.webdriver.common.action_chains import ActionChains
element = sb.find_element("#element")
ActionChains(sb.driver).move_to_element(element).click().perform()
```

### 处理动态 ID

有些网站使用动态生成的 ID，可以使用部分属性或 XPath 定位：

```python
# 使用部分属性匹配
sb.click("button[id*='submit']")  # ID 包含 "submit"

# 使用 XPath 和 contains 函数
sb.click("//button[contains(@id, 'submit')]")

# 使用其他稳定属性
sb.click("button.submit-button")
```

### 处理隐藏元素

有些元素可能被隐藏或覆盖：

```python
# 检查元素是否可见
if sb.is_element_visible("#element"):
    sb.click("#element")
else:
    # 尝试使用 JavaScript 移除覆盖元素
    sb.execute_script("document.querySelector('.overlay').style.display = 'none';")
    sb.click("#element")
```

### 处理验证码

除了使用 CDP 模式外，还可以：

```python
# 检测是否存在验证码
if sb.is_element_present("iframe[title*='recaptcha']"):
    # 方法 1: 使用 UC 模式的自动点击功能
    sb.uc_gui_click_captcha()
    
    # 方法 2: 等待用户手动解决（适用于开发测试）
    print("请手动解决验证码...")
    sb.sleep(30)  # 给用户 30 秒时间解决验证码
```

### 处理网络问题

处理网络超时或连接问题：

```python
from seleniumbase import SB
import time

def retry_on_failure(max_retries=3, delay=5):
    for attempt in range(max_retries):
        try:
            with SB(uc=True) as sb:
                sb.open("https://example.com")
                # 执行其他操作...
                return "成功"
        except Exception as e:
            print(f"尝试 {attempt+1}/{max_retries} 失败: {str(e)}")
            if attempt < max_retries - 1:
                print(f"等待 {delay} 秒后重试...")
                time.sleep(delay)
    
    return "所有尝试均失败"

result = retry_on_failure()
print(result)
```

## 最佳实践

### 爬虫效率优化

1. **使用无头模式**：对于纯数据抓取，使用无头模式可以节省资源。

```python
with SB(uc=True, headless=True) as sb:
    # 爬虫代码...
```

2. **并行抓取**：使用多线程或多进程并行抓取多个页面。

3. **减少不必要的等待**：只在必要时使用 `sleep()`，优先使用显式等待。

4. **使用缓存**：缓存已抓取的数据，避免重复请求。

```python
import os
import json

def get_cached_or_fetch(url, cache_file):
    if os.path.exists(cache_file):
        with open(cache_file, "r", encoding="utf-8") as f:
            return json.load(f)
    
    with SB(uc=True) as sb:
        sb.open(url)
        # 抓取数据...
        data = {...}  # 抓取的数据
        
        # 缓存数据
        with open(cache_file, "w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        
        return data
```

### 反爬虫策略

1. **模拟真实用户行为**：随机等待、随机滚动、不规则点击。

2. **使用代理**：轮换使用多个代理 IP。

```python
proxies = [
    "socks5://proxy1.example.com:1080",
    "socks5://proxy2.example.com:1080",
    "socks5://proxy3.example.com:1080"
]

import random
proxy = random.choice(proxies)

with SB(uc=True, proxy=proxy) as sb:
    # 爬虫代码...
```

3. **设置合理的请求间隔**：避免频繁请求同一网站。

4. **修改浏览器指纹**：使用 CDP 模式修改浏览器特征。

```python
with SB(uc=True) as sb:
    # 修改 User-Agent
    user_agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    sb.execute_cdp_cmd("Network.setUserAgentOverride", {"userAgent": user_agent})
    
    # 修改其他浏览器特征
    sb.execute_script("""
        Object.defineProperty(navigator, 'webdriver', {
            get: () => false
        });
    """)
    
    # 继续爬虫操作...
```

### 错误处理与恢复

1. **使用 try-except 捕获异常**：处理可能的错误情况。

```python
try:
    with SB(uc=True) as sb:
        sb.open("https://example.com")
        # 爬虫代码...
except Exception as e:
    print(f"发生错误: {str(e)}")
    # 错误处理逻辑...
```

2. **实现断点续爬**：记录已处理的页面，支持从中断处继续。

```python
import json
import os

def scrape_with_resume(start_url, max_pages=100):
    # 加载已抓取的页面记录
    resume_file = "resume_state.json"
    if os.path.exists(resume_file):
        with open(resume_file, "r") as f:
            state = json.load(f)
        visited_urls = set(state["visited"])
        results = state["results"]
        next_url = state["next_url"]
    else:
        visited_urls = set()
        results = []
        next_url = start_url
    
    try:
        with SB(uc=True) as sb:
            while next_url and len(visited_urls) < max_pages:
                if next_url in visited_urls:
                    # 避免循环
                    break
                
                print(f"正在抓取: {next_url}")
                sb.open(next_url)
                visited_urls.add(next_url)
                
                # 抓取当前页面数据
                page_data = {
                    "url": next_url,
                    "title": sb.get_title(),
                    "content": sb.get_text("main")
                }
                results.append(page_data)
                
                # 查找下一页链接
                if sb.is_element_present("a.next-page"):
                    next_url = sb.get_attribute("a.next-page", "href")
                else:
                    next_url = None
                
                # 保存当前状态（断点）
                state = {
                    "visited": list(visited_urls),
                    "results": results,
                    "next_url": next_url
                }
                with open(resume_file, "w") as f:
                    json.dump(state, f)
        
        return results
    
    except Exception as e:
        print(f"抓取中断: {str(e)}")
        # 状态已保存，下次可以继续
        return results

# 执行带断点续爬的抓取
results = scrape_with_resume("https://example.com/page/1")
```

通过本教程，您应该已经掌握了 SeleniumBase 的基本用法和高级功能，特别是在爬虫场景下的应用。SeleniumBase 强大的 API 和反爬虫能力使其成为网页自动化和数据抓取的理想工具。随着实践的深入，您将能够应对更复杂的爬虫挑战，并开发出更高效、更可靠的爬虫程序。

from seleniumbase import SB
import random
import time

def scrape_weibo_with_uc():
    with SB(
        uc=True,
        headless=False,
        user_data_dir="./chrome_profile",  # 使用持久化用户数据
        browser="chrome",
        incognito=False,  # 非无痕模式更自然
        agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        proxy="http://user:pass@proxy_ip:port",  # 使用高质量代理
        disable_csp=True,
        do_not_track=True,
        block_images=True,
        undetectable=True,  # 启用反检测
        uc_cdp=True,  # 启用CDP覆盖
        uc_subprocess=True,
    ) as sb:
        # 设置更自然的浏览行为
        sb.set_window_rect(
            x=random.randint(0, 500),
            y=random.randint(0, 500),
            width=random.randint(1200, 1920),
            height=random.randint(800, 1080)
        )
        
        # 随机化鼠标移动
        def random_mouse_move():
            sb.move_to("body", x=random.randint(0, 500), y=random.randint(0, 500))
            sb.sleep(random.uniform(0.5, 2.0))
            
        # 访问前随机操作
        random_mouse_move()
        sb.open_new_window()
        sb.open("https://www.baidu.com")
        sb.type('input[name="wd"]', "test", delay=random.uniform(0.1, 0.5))
        sb.sleep(random.uniform(1, 3))
        sb.go_back()
        
        # 访问微博
        sb.open("https://weibo.com")
        sb.sleep(random.uniform(2, 5))
        
        # 处理可能的验证码
        if "验证码" in sb.get_page_source():
            sb.sleep(10)  # 手动处理验证码
            # 或者使用第三方验证码服务
            
        # 搜索操作
        search_box = 'input[type="text"][node-type="searchInput"]'
        sb.click(search_box)
        sb.sleep(random.uniform(0.5, 1.5))
        sb.type(search_box, "Python爬虫", delay=random.uniform(0.1, 0.3))
        sb.sleep(random.uniform(1, 2))
        sb.click('a[node-type="searchSubmit"]')
        
        # 高级滚动策略
        last_height = sb.execute_script("return document.body.scrollHeight")
        scroll_attempts = 0
        while scroll_attempts < 10:
            # 随机滚动距离
            scroll_distance = random.randint(300, 800)
            sb.execute_script(f"window.scrollBy(0, {scroll_distance})")
            
            # 随机等待时间
            wait_time = random.uniform(0.5, 3.0)
            sb.sleep(wait_time)
            
            # 随机点击页面空白处
            if random.random() > 0.7:
                sb.click('body', x=random.randint(0, 500), y=random.randint(0, 300))
                
            # 检查是否到达底部
            new_height = sb.execute_script("return document.body.scrollHeight")
            if new_height == last_height:
                scroll_attempts += 1
            else:
                scroll_attempts = 0
            last_height = new_height
            
            # 随机移动鼠标
            if random.random() > 0.5:
                random_mouse_move()
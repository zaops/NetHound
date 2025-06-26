import os
import time
from pathlib import Path
from collections import defaultdict

def get_dir_size(path):
    total = 0
    try:
        # 遍历目录下的所有文件和文件夹
        for entry in os.scandir(path):
            try:
                if entry.is_file():
                    total += entry.stat().st_size
                elif entry.is_dir():
                    total += get_dir_size(entry.path)
            except (PermissionError, FileNotFoundError):
                continue
    except (PermissionError, FileNotFoundError):
        return 0
    return total

def format_size(size):
    # 转换字节到合适的单位
    for unit in ['B', 'KB', 'MB', 'GB']:
        if size < 1024:
            return f"{size:.2f} {unit}"
        size /= 1024
    return f"{size:.2f} TB"

def get_clean_suggestion(folder_name, size_gb):
    # 根据文件夹名称和大小给出清理建议
    suggestions = {
        'Temp': '建议清理，临时文件可以安全删除',
        'Cache': '建议定期清理，缓存文件可以安全删除',
        'Local': '需谨慎清理，包含应用程序数据',
        'LocalLow': '可以适当清理，一般存放低权限应用数据',
        'Roaming': '需谨慎清理，包含重要的用户配置',
        'Microsoft': '需谨慎清理，包含系统组件数据',
        'Packages': '需谨慎清理，包含应用商店应用数据',
        'Programs': '需谨慎清理，包含已安装程序',
        'Downloads': '可以清理，注意备份重要文件',
        'Logs': '可以清理，日志文件',
        'Updates': '可以清理，但建议保留最近的更新',
        'Temp': '建议清理，临时文件',
        'History': '可以清理，历史记录文件'
    }
    
    for key in suggestions:
        if key.lower() in folder_name.lower():
            return suggestions[key]
    
    if size_gb > 10:
        return '建议检查具体内容，体积较大'
    elif size_gb > 1:
        return '可以考虑检查内容'
    else:
        return '体积较小，可以保留'

def analyze_directory(path, min_size_mb=100):
    results = []
    try:
        for entry in os.scandir(path):
            if entry.is_dir():
                try:
                    size = get_dir_size(entry.path)
                    size_gb = size / (1024**3)
                    size_mb = size / (1024**2)
                    
                    # 只记录大于指定大小的文件夹
                    if size_mb >= min_size_mb:
                        results.append({
                            'name': entry.name,
                            'path': entry.path,
                            'size': size,
                            'size_formatted': format_size(size),
                            'suggestion': get_clean_suggestion(entry.name, size_gb)
                        })
                        
                        # 递归分析子文件夹
                        sub_results = analyze_directory(entry.path, min_size_mb)
                        if sub_results:
                            results.extend(sub_results)
                except (PermissionError, FileNotFoundError):
                    continue
    except (PermissionError, FileNotFoundError):
        return []
    return results

def main():
    appdata_path = os.path.expandvars('%LOCALAPPDATA%')
    parent_path = os.path.dirname(appdata_path)  # 获取包含Local、Roaming等的父目录
    
    print(f"\n开始分析目录: {parent_path}\n")
    print("-" * 80)
    
    # 分析目录，只显示大于100MB的文件夹
    folders_info = analyze_directory(parent_path, min_size_mb=100)
    
    # 按大小排序
    folders_info.sort(key=lambda x: x['size'], reverse=True)
    
    # 打印结果
    print(f"{'文件夹':<50} {'大小':<15} {'清理建议'}")
    print("-" * 80)
    for folder in folders_info:
        # 缩进显示子文件夹
        indent = ' ' * (folder['path'].count(os.sep) - parent_path.count(os.sep))
        name_display = indent + folder['name']
        print(f"{name_display:<50} {folder['size_formatted']:<15} {folder['suggestion']}")

if __name__ == '__main__':
    main()
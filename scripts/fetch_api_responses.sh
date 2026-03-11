#!/bin/bash
# 使用 curl 调用百炼 API，将单轮、多轮、流式响应保存到 data 目录
# 使用前请设置: export DASHSCOPE_API_KEY=sk-xxx DASHSCOPE_APP_ID=your-app-id

set -e
export LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8  # 确保 curl 输出 UTF-8
cd "$(dirname "$0")/.."
DATA_DIR="data"
API_KEY="${DASHSCOPE_API_KEY}"
APP_ID="${DASHSCOPE_APP_ID}"
BASE_URL="https://dashscope.aliyuncs.com/api/v1/apps/${APP_ID}/completion"

if [ -z "$API_KEY" ] || [ -z "$APP_ID" ]; then
  echo "请设置环境变量: DASHSCOPE_API_KEY, DASHSCOPE_APP_ID"
  exit 1
fi

mkdir -p "$DATA_DIR"

echo "========== 1. 单轮对话 =========="
curl -s -X POST "$BASE_URL" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "input": {"prompt": "考勤异常怎么办？"},
    "parameters": {"has_thoughts": true},
    "debug": {}
  }' | tee "$DATA_DIR/single_turn.json" | jq . 2>/dev/null || cat "$DATA_DIR/single_turn.json"
echo "已保存到 $DATA_DIR/single_turn.json"

echo ""
echo "========== 2. 多轮对话（第一轮）=========="
curl -s -X POST "$BASE_URL" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "input": {"prompt": "精臣科技的开票信息"},
    "parameters": {"has_thoughts": true},
    "debug": {}
  }' > "$DATA_DIR/multi_turn_1.json"
SESSION_ID=$(jq -r '.output.session_id // empty' "$DATA_DIR/multi_turn_1.json")
echo "session_id: $SESSION_ID"
jq . "$DATA_DIR/multi_turn_1.json" 2>/dev/null || cat "$DATA_DIR/multi_turn_1.json"

echo "已保存到 $DATA_DIR/multi_turn_1.json"
if [ -n "$SESSION_ID" ]; then
  echo ""
  echo "========== 2. 多轮对话（第二轮）=========="
  curl -s -X POST "$BASE_URL" \
    -H "Authorization: Bearer $API_KEY" \
    -H "Content-Type: application/json" \
    -d "{
      \"input\": {\"prompt\": \"精臣智慧的呢？\", \"session_id\": \"$SESSION_ID\"},
      \"parameters\": {\"has_thoughts\": true},
      \"debug\": {}
    }" > "$DATA_DIR/multi_turn_2.json"
  jq . "$DATA_DIR/multi_turn_2.json" 2>/dev/null || cat "$DATA_DIR/multi_turn_2.json"
  echo "已保存到 $DATA_DIR/multi_turn_2.json"
fi

echo ""
echo "========== 3. 流式对话 =========="
curl -s -X POST "$BASE_URL" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -H "X-DashScope-SSE: enable" \
  -d '{
    "input": {"prompt": "巴黎出差的差旅费是多少？"},
    "parameters": {"incremental_output": true, "has_thoughts": true},
    "debug": {}
  }' > "$DATA_DIR/stream.txt"
echo "已保存到 $DATA_DIR/stream.txt"
echo "流式响应行数: $(wc -l < "$DATA_DIR/stream.txt")"

echo ""
echo "========== 完成 =========="
ls -la "$DATA_DIR"/

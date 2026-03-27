#!/bin/sh
input=$(cat)

# Parse JSON fields
model=$(echo "$input" | jq -r '.model.display_name // "---"')
model_id=$(echo "$input" | jq -r '.model.id // ""')
ctx_used=$(echo "$input" | jq -r '.context_window.used_percentage // 0')
ctx_remain=$(echo "$input" | jq -r '.context_window.remaining_percentage // 100')
ctx_size=$(echo "$input" | jq -r '.context_window.context_window_size // 0')
input_tok=$(echo "$input" | jq -r '.context_window.total_input_tokens // 0')
output_tok=$(echo "$input" | jq -r '.context_window.total_output_tokens // 0')
cost=$(echo "$input" | jq -r '.cost.total_cost_usd // 0')
duration_ms=$(echo "$input" | jq -r '.cost.total_duration_ms // 0')
lines_add=$(echo "$input" | jq -r '.cost.total_lines_added // 0')
lines_rm=$(echo "$input" | jq -r '.cost.total_lines_removed // 0')
cwd=$(echo "$input" | jq -r '.workspace.current_dir // .cwd // "~"')
version=$(echo "$input" | jq -r '.version // "?"')
exceeds_200k=$(echo "$input" | jq -r '.exceeds_200k // false')
effort=$(jq -r '.effortLevel // "default"' ~/.claude/settings.json 2>/dev/null)

# Git branch - try project_dir first, then cwd
proj_dir=$(echo "$input" | jq -r '.workspace.project_dir // ""')
branch=""
for d in "$proj_dir" "$cwd"; do
  [ -z "$d" ] && continue
  branch=$(git -C "$d" rev-parse --abbrev-ref HEAD 2>/dev/null)
  [ -n "$branch" ] && break
done
[ -z "$branch" ] || [ "$branch" = "HEAD" ] && branch="detached"

# Format duration
duration_s=$((duration_ms / 1000))
if [ "$duration_s" -ge 3600 ]; then
  dur_fmt="$(( duration_s / 3600 ))h$(( (duration_s % 3600) / 60 ))m"
elif [ "$duration_s" -ge 60 ]; then
  dur_fmt="$(( duration_s / 60 ))m$(( duration_s % 60 ))s"
else
  dur_fmt="${duration_s}s"
fi

# Format tokens (k)
fmt_tok() {
  if [ "$1" -ge 1000000 ]; then
    printf "%.1fM" "$(echo "$1 / 1000000" | bc -l)"
  elif [ "$1" -ge 1000 ]; then
    printf "%.1fk" "$(echo "$1 / 1000" | bc -l)"
  else
    printf "%d" "$1"
  fi
}
in_fmt=$(fmt_tok "$input_tok")
out_fmt=$(fmt_tok "$output_tok")

# Format context window size
if [ "$ctx_size" -ge 1000000 ]; then
  ctx_size_fmt="1M"
else
  ctx_size_fmt="200k"
fi

# Colors - hacker green palette
G='\033[38;5;46m'    # bright green
DG='\033[38;5;22m'   # dark green
MG='\033[38;5;34m'   # mid green
LG='\033[38;5;119m'  # lime green
Y='\033[38;5;226m'   # yellow warning
R='\033[38;5;196m'   # red alert
DIM='\033[38;5;240m' # dim gray
W='\033[38;5;255m'   # white
RST='\033[0m'

# Context bar - 20 chars wide
bar_width=20
filled=$(( ctx_used * bar_width / 100 ))
[ "$filled" -gt "$bar_width" ] && filled=$bar_width
empty=$(( bar_width - filled ))

# Color based on usage
if [ "$ctx_used" -ge 80 ]; then
  BAR_COLOR="$R"
  CTX_ICON="!!"
elif [ "$ctx_used" -ge 50 ]; then
  BAR_COLOR="$Y"
  CTX_ICON=">>"
else
  BAR_COLOR="$G"
  CTX_ICON="::"
fi

# Build the bar
bar=""
i=0; while [ "$i" -lt "$filled" ]; do bar="${bar}#"; i=$((i+1)); done
i=0; while [ "$i" -lt "$empty" ]; do bar="${bar}-"; i=$((i+1)); done

# Cost color
cost_int=$(printf "%.0f" "$(echo "$cost * 100" | bc -l)")
if [ "$cost_int" -ge 500 ]; then
  COST_C="$R"
elif [ "$cost_int" -ge 100 ]; then
  COST_C="$Y"
else
  COST_C="$G"
fi

# Effort color
case "$effort" in
  high) EFFORT_C="$R"; effort_tag="HIGH" ;;
  low)  EFFORT_C="$DG"; effort_tag="LOW" ;;
  *)    EFFORT_C="$MG"; effort_tag="MID" ;;
esac

# Line 1: model + effort + context bar
printf "${DIM}[${MG}%s${DIM}|${G}%s${DIM}] ${EFFORT_C}%s ${DIM}ctx${BAR_COLOR}[%s] ${BAR_COLOR}%s%%${RST}\n" \
  "$model" "$ctx_size_fmt" "$effort_tag" "$bar" "$ctx_used"

# Line 2: cost + tokens (spaced) + lines changed + session time + git
printf "${COST_C}\$%.4f ${DIM}in:${LG}%s ${DIM}out:${LG}%s ${DIM}lines:${G}+%s${DIM}/${R}-%s ${DIM}time:${MG}%s ${DIM}git:${DG}(%s)${RST}" \
  "$cost" "$in_fmt" "$out_fmt" "$lines_add" "$lines_rm" "$dur_fmt" "$branch"

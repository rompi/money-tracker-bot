#!/bin/bash

# Git Workflow Script for Money Tracker Bot
# This script enforces the development standards outlined in README.md

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "Not in a git repository"
        exit 1
    fi
}

# Function to run pre-commit checks
run_pre_commit_checks() {
    print_status "Running pre-commit checks..."

    # Check for staged files
    if ! git diff --cached --quiet; then
        print_status "Found staged files, running checks..."
    else
        print_warning "No staged files found. Stage your changes with 'git add' first."
        exit 1
    fi

    # Run tests with coverage
    print_status "Running tests with coverage..."
    if ! make test; then
        print_error "Tests failed. Please fix failing tests before committing."
        print_status "Run 'make test' to see detailed test output."
        exit 1
    fi
    print_success "All tests passed"

    # Format code
    print_status "Formatting code..."
    make fmt
    print_success "Code formatted"

    # Run linting
    print_status "Running linting checks..."
    if ! make lint; then
        print_error "Linting failed. Please fix linting issues before committing."
        exit 1
    fi
    print_success "Linting passed"

    # Check if formatting created any changes
    if ! git diff --quiet; then
        print_warning "Code formatting made changes. Please review and stage them:"
        git status --porcelain
        print_status "Run 'git add .' to stage formatting changes, then run this script again."
        exit 1
    fi
}

# Function to setup git hooks
setup_git_hooks() {
    print_status "Setting up git pre-commit hook..."

    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
echo "🔍 Running pre-commit checks..."

# Run tests
if ! make test; then
    echo "❌ Tests failed. Commit blocked."
    echo "💡 Fix the failing tests and try again."
    exit 1
fi

# Run formatting and linting
make fmt
if ! make lint; then
    echo "❌ Linting failed. Commit blocked."
    echo "💡 Fix the linting issues and try again."
    exit 1
fi

# Check if formatting made changes
if ! git diff --quiet; then
    echo "⚠️  Code formatting made changes. Please stage them:"
    git status --porcelain
    echo "💡 Run 'git add .' and commit again."
    exit 1
fi

echo "✅ All checks passed. Proceeding with commit."
EOF

    chmod +x .git/hooks/pre-commit
    print_success "Git pre-commit hook installed"
}

# Function to create a new feature branch
create_feature_branch() {
    local branch_name="$1"

    if [ -z "$branch_name" ]; then
        print_error "Branch name is required"
        echo "Usage: $0 branch <feature-name>"
        exit 1
    fi

    # Ensure we're on main/master branch
    current_branch=$(git branch --show-current)
    if [[ "$current_branch" != "main" && "$current_branch" != "master" ]]; then
        print_warning "You're not on main/master branch. Current branch: $current_branch"
        read -p "Do you want to continue? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi

    # Create and checkout feature branch
    git checkout -b "feature/$branch_name"
    print_success "Created and switched to branch: feature/$branch_name"
}

# Function to commit with checks
commit_with_checks() {
    local commit_message="$1"

    if [ -z "$commit_message" ]; then
        print_error "Commit message is required"
        echo "Usage: $0 commit \"Your commit message\""
        exit 1
    fi

    # Run pre-commit checks
    run_pre_commit_checks

    # Commit the changes
    git commit -m "$commit_message"
    print_success "Changes committed successfully"

    # Show current branch and suggest next steps
    current_branch=$(git branch --show-current)
    print_status "Current branch: $current_branch"
    if [[ "$current_branch" == feature/* ]]; then
        print_status "💡 Next steps:"
        echo "   1. Push branch: git push origin $current_branch"
        echo "   2. Create pull request on GitHub"
    fi
}

# Function to show usage
show_usage() {
    echo "Git Workflow Script for Money Tracker Bot"
    echo ""
    echo "Usage: $0 <command> [arguments]"
    echo ""
    echo "Commands:"
    echo "  setup                    - Install git pre-commit hooks"
    echo "  check                    - Run pre-commit checks without committing"
    echo "  branch <name>            - Create new feature branch"
    echo "  commit \"<message>\"       - Run checks and commit with message"
    echo "  help                     - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 setup"
    echo "  $0 branch transaction-categories"
    echo "  $0 check"
    echo "  $0 commit \"feat: add transaction validation\""
    echo ""
    echo "Quality Requirements:"
    echo "  ✅ Tests must pass with ≥85% coverage"
    echo "  ✅ Code must be formatted (make fmt)"
    echo "  ✅ Linting must pass (make lint)"
    echo "  ✅ Package AI.md files must be updated"
    echo "  ✅ Function docstrings with examples required"
}

# Main script logic
main() {
    check_git_repo

    case "${1:-help}" in
        "setup")
            setup_git_hooks
            ;;
        "check")
            run_pre_commit_checks
            print_success "All pre-commit checks passed!"
            ;;
        "branch")
            create_feature_branch "$2"
            ;;
        "commit")
            commit_with_checks "$2"
            ;;
        "help"|"--help"|"-h")
            show_usage
            ;;
        *)
            print_error "Unknown command: $1"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
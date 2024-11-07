import subprocess
import requests
import os
import sys
import logging
from pathlib import Path
from requests.exceptions import RequestException
from typing import List, Optional

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

class TestGenerator:
    def __init__(self):
        self.api_key = os.getenv('OPENAI_API_KEY')
        self.model = os.getenv('OPENAI_MODEL', 'gpt-4-turbo-preview')
        
        try:
            self.max_tokens = int(os.getenv('OPENAI_MAX_TOKENS', '2000'))
        except ValueError:
            logging.error("Invalid value for OPENAI_MAX_TOKENS. Using default value: 2000")
            self.max_tokens = 2000

        if not self.api_key:
            raise ValueError("OPENAI_API_KEY environment variable is not set")

    def get_changed_files(self) -> List[str]:
        """Retrieve list of changed files passed as command-line arguments."""
        if len(sys.argv) <= 1:
            return []
        return [f.strip() for f in sys.argv[1:] if f.strip()]

    def detect_language(self, file_name: str) -> str:
        """Detect programming language based on file extension."""
        extensions = {
            '.py': 'Python',
            '.js': 'JavaScript',
            '.ts': 'TypeScript',
            '.java': 'Java',
            '.cpp': 'C++',
            '.cs': 'C#',
            '.go': 'Go'
        }
        _, ext = os.path.splitext(file_name)
        return extensions.get(ext.lower(), 'Unknown')

    def get_test_framework(self, language: str) -> str:
        """Get the appropriate test framework based on language."""
        frameworks = {
            'Python': 'pytest',
            'JavaScript': 'jest',
            'TypeScript': 'jest',
            'Java': 'JUnit',
            'C++': 'Google Test',
            'C#': 'NUnit',
            'Go': 'testing'
        }
        return frameworks.get(language, 'unknown')
    
    def get_related_files(self, language: str, file_name: str) -> List[str]:
        """Identify related files based on import statements or includes."""
        related_files = []
        
        try:
            if language in ["Python", "JavaScript", "TypeScript"]:
                with open(file_name, 'r') as f:
                    for line in f:
                        if 'import ' in line or 'from ' in line or 'require(' in line:
                            parts = line.split()
                            for part in parts:
                                # Check for relative imports
                                if len(part) > 1 and part.startswith(".") and not part.startswith(".."):
                                    path = part.replace(".", "")
                                    for ext in ('.py', '.js', '.ts'):
                                        potential_file = f"{path}{ext}"
                                        if Path(potential_file).exists():
                                            related_files.append(potential_file)
                                            break
                                elif '.' in part:
                                    path = part.replace(".", "/")
                                    for ext in ('.py', '.js', '.ts'):
                                        potential_file = f"{path}{ext}"
                                        if Path(potential_file).exists():
                                            related_files.append(potential_file)
                                            break
                                else:
                                    if part.endswith(('.py', '.js', '.ts')) and Path(part).exists():
                                        related_files.append(part)
                                    elif part.isidentifier():
                                        base_name = part.lower()
                                        for ext in ('.py', '.js', '.ts'):
                                            potential_file = f"{base_name}{ext}"
                                            if Path(potential_file).exists():
                                                related_files.append(potential_file)
                                                break
            elif language in ["C++", "C#"]:
                # Placeholder for C++ and C# related files logic
                pass

        except Exception as e:
            logging.error(f"Error identifying related files in {file_name}: {e}")
        
        return related_files

    def get_related_test_files(self, language: str, file_name: str) -> List[str]:
        """Identify related test files based on import statements or includes."""
        related_test_files = []
        try:
            if language == "Python":
                directory = Path(os.path.dirname(os.path.abspath(file_name)))
                test_files = list(directory.rglob("tests.py")) + list(directory.rglob("test.py")) + \
                             list(directory.rglob("test_*.py")) + list(directory.rglob("*_test.py"))
                for file in test_files:
                    with open(file, 'r') as f:
                        for line in f:
                            if 'from ' in line or 'import ' in line or 'require(' in line:
                                parts = line.split()
                                for part in parts:
                                    if len(part) > 1 and part.startswith(".") and not part.startswith(".."):
                                        path = part.replace(".", "")
                                        for ext in ('.py', '.js', '.ts'):
                                            potential_file = f"{path}{ext}"
                                            if Path(potential_file).exists() and (Path(file_name).name in potential_file):
                                                related_test_files.append(str(file))
                                                break
                                    elif '.' in part:
                                        path = part.replace(".", "/")
                                        for ext in ('.py', '.js', '.ts'):
                                            potential_file = f"{path}{ext}"
                                            if Path(potential_file).exists() and (Path(file_name).name in potential_file):
                                                related_test_files.append(str(file))
                                                break
                                    else:
                                        if part.endswith(('.py', '.js', '.ts')) and Path(part).exists() and (Path(file_name).name in part):
                                            related_test_files.append(str(file))
                                        elif part.isidentifier():
                                            base_name = part.lower()
                                            for ext in ('.py', '.js', '.ts', '.js'):
                                                potential_file = f"{base_name}{ext}"
                                                if Path(potential_file).exists() and (Path(file_name).name in potential_file):
                                                    related_test_files.append(str(file))
                                                    break
            # Add other language test file identification as needed
        except Exception as e:
            logging.error(f"Error identifying related test files in {file_name}: {e}")
        
        # Limit to 1 related test file to prevent excessive processing
        limited_test_files = related_test_files[:1]
        return limited_test_files

    def get_package_name(self, file_path: Path) -> str:
        """Extract the package name from a Go source file."""
        try:
            with open(file_path, 'r') as f:
                for line in f:
                    line = line.strip()
                    if line.startswith("package "):
                        return line.split()[1]
        except Exception as e:
            logging.error(f"Error reading package name from {file_path}: {e}")
        return "main"  # Default package name if not found

    def generate_coverage_report_for_go(self):
        """Generate a Go coverage report after processing all Go test files."""
        repo_root = Path.cwd()
        logging.info(f"Repository root is: {repo_root}")
        go_mod = repo_root / "go.mod"
        if not go_mod.exists():
            logging.error(f"'go.mod' not found in repository root: {repo_root}")
            return

        report_out = repo_root / "coverage_report.out"
        report_html = repo_root / "coverage_report.html"

        try:
            logging.info("Running 'go test' in repository root...")
            subprocess.run(
                ["go", "test", "./...", "-coverprofile", str(report_out)],
                cwd=repo_root,
                check=True
            )
            logging.info(f"Generated cover profile at: {report_out}")

            # Convert coverprofile to human-readable HTML format
            subprocess.run(
                ["go", "tool", "cover", "-html", str(report_out), "-o", str(report_html)],
                check=True
            )
            logging.info(f"HTML coverage report generated at {report_html}")

            # Verify coverage report files
            if report_out.exists():
                logging.info(f"Cover profile file exists with size {report_out.stat().st_size} bytes.")
            else:
                logging.error(f"Cover profile file {report_out} does not exist.")

            if report_html.exists():
                logging.info(f"HTML coverage report exists with size {report_html.stat().st_size} bytes.")
            else:
                logging.error(f"HTML coverage report {report_html} does not exist.")

        except subprocess.CalledProcessError as e:
            logging.error(f"Error generating Go coverage report: {e}")

    def generate_coverage_report(self, test_file: Path, language: str):
        """Generate a code coverage report and save it as a text or HTML file."""
        # Confirm repository root
        repo_root = Path.cwd()
        logging.info(f"Repository root is: {repo_root}")
        go_mod = repo_root / "go.mod"
        if not go_mod.exists():
            logging.error(f"'go.mod' not found in repository root: {repo_root}")
            return

        report_out = repo_root / "coverage_report.out"
        report_html = repo_root / "coverage_report.html"
        
        try:
            # Run tests with coverage based on language
            if language == "Python":
                subprocess.run(
                    ["coverage", "run", str(test_file)],
                    check=True
                )
                subprocess.run(
                    ["coverage", "report", "-m", "--omit=*/site-packages/*"],
                    stdout=open(repo_root / "coverage_report.txt", "w"),
                    check=True
                )
            elif language == "JavaScript":
                subprocess.run(
                    ["npx", "jest", "--coverage", "--config=path/to/jest.config.js"],
                    stdout=open(repo_root / "coverage_report.txt", "w"),
                    check=True
                )
            elif language == "Go":
                # Delegate to generate_coverage_report_for_go
                self.generate_coverage_report_for_go()
                return

            # Add additional commands for other languages here
            logging.info(f"Code coverage report saved to {report_html if language == 'Go' else report_out}")
        
        except subprocess.CalledProcessError as e:
            logging.error(f"Error generating coverage report for {test_file}: {e}")

    def ensure_coverage_installed(self, language: str):
        """
        Ensures that the appropriate coverage tool for the given programming language is installed.
        Logs messages for each step.
        """
        try:
            if language.lower() == 'python':
                # Check if 'coverage' is installed for Python
                subprocess.check_call([sys.executable, '-m', 'pip', 'show', 'coverage'])
                logging.info(f"Coverage tool for Python is already installed.")
            elif language.lower() == 'javascript':
                # Check if 'jest' coverage is available for JavaScript
                subprocess.check_call(['npx', 'jest', '--version'])
                logging.info(f"Coverage tool for JavaScript (jest) is already installed.")
            elif language.lower() == 'java':
                # Check if 'jacoco' is available for Java (typically part of the build process)
                logging.info("Make sure Jacoco is configured in your Maven/Gradle build.")
                # Optionally you can add a check for specific build tool (Maven/Gradle) commands here
            elif language.lower() == 'ruby':
                # Check if 'simplecov' is installed for Ruby
                subprocess.check_call(['gem', 'list', 'simplecov'])
                logging.info(f"Coverage tool for Ruby (simplecov) is already installed.")
            elif language.lower() == 'go':
                # Go has built-in coverage tools; no additional installation needed
                logging.info("Go's built-in coverage tools are available.")
            else:
                logging.warning(f"Coverage tool check is not configured for {language}. Please add it manually.")
                return

        except subprocess.CalledProcessError:
            logging.error(f"Coverage tool for {language} is not installed or not accessible.")

            try:
                if language.lower() == 'python':
                    subprocess.check_call([sys.executable, '-m', 'pip', 'install', 'coverage'])
                    logging.info(f"Coverage tool for Python has been installed.")
                elif language.lower() == 'javascript':
                    subprocess.check_call(['npm', 'install', 'jest'])
                    logging.info(f"Coverage tool for JavaScript (jest) has been installed.")
                elif language.lower() == 'ruby':
                    subprocess.check_call(['gem', 'install', 'simplecov'])
                    logging.info(f"Coverage tool for Ruby (simplecov) has been installed.")
                else:
                    logging.error(f"Could not install coverage tool for {language} automatically. Please install manually.")
            except subprocess.CalledProcessError:
                logging.error(f"Failed to install the coverage tool for {language}. Please install it manually.")

    def create_prompt(self, file_name: str, language: str) -> Optional[str]:
        """Create a language-specific prompt for test generation with accurate module and import names in related content."""
        try:
            with open(file_name, 'r') as f:
                code_content = f.read()
        except Exception as e:
            logging.error(f"Error reading file {file_name}: {e}")
            return None

        # Extract package name for Go
        package_name = ""
        if language == "Go":
            package_name = self.get_package_name(Path(file_name))

        # Gather related files and embed imports in each file's content
        related_files = self.get_related_files(language, file_name)
        related_content = ""

        # Log related files to confirm detection
        if related_files:
            logging.info(f"Related files for {file_name}: {related_files}")
        else:
            logging.info(f"No related files found for {file_name} to reference")
        for related_file in related_files:
            try:
                with open(related_file, 'r') as rf:
                    file_content = rf.read()
                    
                    # Generate the correct module path for import statements
                    module_path = str(Path(related_file).with_suffix('')).replace('/', '.')
                    if language == 'Go':
                        import_statement = f'import "{module_path}"'
                    else:
                        import_statement = f"import {module_path}"
                    
                    # Append file content with embedded import statement
                    related_content += f"\n\n// Module: {module_path}\n{import_statement}\n{file_content}"
                    logging.info(f"Included content from related file: {related_file} as module {module_path}")
            except Exception as e:
                logging.error(f"Error reading related file {related_file}: {e}")

        # Gather additional context from related test files
        
        related_test_files = self.get_related_test_files(language, file_name)
        related_test_content = ""
        # Log related files to confirm detection
        if related_test_files:
            logging.info(f"Related Test files for {file_name}: {related_test_files}")
        else:
            logging.info(f"No related test files found for {file_name} to reference")
        for related_test_file in related_test_files:
            try:
                with open(related_test_file, 'r') as rf:
                    file_content = rf.read()
                    related_test_content += f"\n\n// Related test file: {related_test_file}\n{file_content}"
                    logging.info(f"Included content from related test file: {related_test_file}")
            except Exception as e:
                logging.error(f"Error reading related test file {related_test_file}: {e}")

        # Add the file name and package name at the top of the prompt for Go
        framework = self.get_test_framework(language)
        if language == "Go" and package_name:
            package_info = f"The test code must start with the correct package declaration: package {package_name}."
        else:
            package_info = "Ensure the test code starts with the correct package declaration."

        prompt = f"""Generate comprehensive unit tests for the following {language} file: {file_name} using {framework}.

Requirements:
1. Include edge cases, normal cases, and error cases.
2. Use mocking where appropriate for external dependencies.
3. Include setup and teardown if needed.
4. Add descriptive test names and docstrings.
5. Follow {framework} best practices.
6. Ensure high code coverage.
7. Test both success and failure scenarios.
8. {package_info}

Code to test (File: {file_name}):

{code_content}

Related context:

{related_content}

Related test cases:
{related_test_content}

Generate only the test code without any explanations or notes. Ensure that the test file includes at least one valid test function."""

        logging.info(f"Created prompt for {file_name} with length {len(prompt)} characters")
        return prompt

    def call_openai_api(self, prompt: str) -> Optional[str]:
        """Call OpenAI API to generate test cases."""
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}'
        }
        
        data = {
            'model': self.model,
            'messages': [
                {
                    "role": "system",
                    "content": "You are a senior software engineer specialized in writing comprehensive test suites."
                },
                {
                    "role": "user",
                    "content": prompt
                }
            ],
            'max_tokens': self.max_tokens,
            'temperature': 0.7
        }

        try:
            response = requests.post(
                'https://api.openai.com/v1/chat/completions',
                headers=headers,
                json=data,
                timeout=60
            )
            response.raise_for_status()
            generated_text = response.json()['choices'][0]['message']['content']
            normalized_text = generated_text.replace('“', '"').replace('”', '"').replace("‘", "'").replace("’", "'")
            if normalized_text.startswith('```'):
                first_newline_index = normalized_text.find('\n', 3)
                if first_newline_index != -1:
                    normalized_text = normalized_text[first_newline_index+1:]
                else:
                    normalized_text = normalized_text[3:]
                if normalized_text.endswith('```'):
                    normalized_text = normalized_text[:-3]
            return normalized_text.strip()
        except RequestException as e:
            logging.error(f"API request failed: {e}")
            return None

    def save_test_cases(self, file_name: str, test_cases: str, language: str) -> Path:
        """Save generated test cases to the same directory as the source file."""
        source_path = Path(file_name)
        source_dir = source_path.parent
        base_name = source_path.stem
        if language == 'Go':
            if not base_name.endswith("_test"):
                base_name = f"{base_name}_test"  # Ensure the test file ends with _test for Go
            extension = '.go'
        else:
            if not base_name.endswith("_test"):
                base_name = f"test_{base_name}"
            extension = source_path.suffix
        test_file = source_dir / f"{base_name}{extension}"

        try:
            with open(test_file, 'w', encoding='utf-8') as f:
                f.write(test_cases)
            logging.info(f"Test cases saved to {test_file}")
        except Exception as e:
            logging.error(f"Error saving test cases to {test_file}: {e}")

        if test_file.exists():
            logging.info(f"File {test_file} exists with size {test_file.stat().st_size} bytes.")
            # Read first few lines
            try:
                with open(test_file, 'r') as f:
                    first_lines = ''.join([next(f) for _ in range(5)])
                logging.info(f"First lines of {test_file}:\n{first_lines}")
            except Exception as e:
                logging.error(f"Error reading first lines of {test_file}: {e}")
        else:
            logging.error(f"File {test_file} was not created.")
        return test_file

    def run(self):
        """Main execution method."""
        changed_files = self.get_changed_files()
        if not changed_files:
            logging.info("No files changed.")
            return

        # Flag to determine if Go coverage report needs to be generated
        go_tests_needed = False

        for file_name in changed_files:
            if file_name != "generate_tests.py":
                try:
                    language = self.detect_language(file_name)
                    if language == 'Unknown':
                        logging.warning(f"Unsupported file type: {file_name}")
                        continue

                    logging.info(f"Processing {file_name} ({language})")
                    prompt = self.create_prompt(file_name, language)

                    if prompt:
                        test_cases = self.call_openai_api(prompt)

                        if test_cases:
                            test_cases = test_cases.replace("“", '"').replace("”", '"')
                            test_file = self.save_test_cases(file_name, test_cases, language)

                            self.ensure_coverage_installed(language)

                            if language == 'Go':
                                go_tests_needed = True
                            else:
                                self.generate_coverage_report(test_file, language)
                        else:
                            logging.error(f"Failed to generate test cases for {file_name}")
                except Exception as e:
                    logging.error(f"Error processing {file_name}: {e}")

        # Handle Go coverage after processing all files
        if go_tests_needed:
            self.generate_coverage_report_for_go()

if __name__ == '__main__':
    try:
        generator = TestGenerator()
        generator.run()
    except Exception as e:
        logging.error(f"Fatal error: {e}")
        sys.exit(1)

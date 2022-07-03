import os.path
import pathlib
from typing import Any

from jinja2 import Environment
from plumbum import cli
from plumbum.cmd import rm  # type: ignore

from config import *

log = logging.getLogger()


class Manager(cli.Application):
    """Aiges CI Manager"""

    PROGNAME: str = "build.py"
    VERSION: str = "0.0.1"

    manifest = {}
    ci = None

    def main(self):
        if not self.nested_command:  # will be ``None`` if no sub-command follows
            log.fatal("No subcommand given!")
            print()
            self.help()
            return 1
        elif len(self.nested_command[1]) < 2 and any(
                "generate" in arg for arg in self.nested_command[1]
        ):
            log.error(
                "Subcommand 'generate' missing  required arguments! use 'generate --help'"
            )
            return 1


class ImageTag(object):
    def __init__(self, cuda, python="3.9.13", golang="1.17", distro="ubuntu1804"):
        self.cuda = cuda
        self.python = python
        self.golang = golang
        self.distro = distro

    def __str__(self):
        return "{cuda}-{golang}-{python}-{distro}".format(cuda=self.cuda, golang=self.golang, python=self.python,
                                                          distro=self.distro)


@Manager.subcommand("generate")  # type: ignore
class ManagerGenerate(Manager):
    DESCRIPTION = "Generate Dockerfiles from templates."

    parent: Manager

    vars = {}

    matrix = []
    template_env: Any = Environment(
        extensions=["jinja2.ext.do", "jinja2.ext.loopcontrols"],
        trim_blocks=True,
        lstrip_blocks=True,
    )
    template: Any
    generate_all: Any = cli.Flag(
        ["--all"],
        help="Generate all of the templates.",
    )

    use_github: Any = cli.Flag(
        ["--use_github"],
        help="If Using Github Actions",
    )

    distro: Any = cli.SwitchAttr(
        "--os-name",
        str,
        group="Targeted",
        excludes=["--all", ],
        help="The distro to use.",
        default=None,
    )

    distro_version: Any = cli.SwitchAttr(
        "--os-version",
        str,
        group="Targeted",
        excludes=["--all"],
        help="The distro version",
        default=None,
    )

    cuda_version: Any = cli.SwitchAttr(
        "--cuda-version",
        str,
        excludes=["--all"],
        group="Targeted",
        help="The cuda version to use. Example: '11.2'",
        default=None,
    )

    def matched(self, key):
        match = self.cuda_version_regex.match(key)
        if match:
            return match

    # extracts arbitrary keys and inserts them into the templating context
    def extract_keys(self, val, arch=None):
        pass

    # For cudnn templates, we need a custom template context
    def output_cudnn_template(self, cudnn_version_name, input_template, output_path):
        pass

    def prepare_context(self):

        # The templating context. This data structure is used to fill the templates.
        self.vars = {
            "registry": self.get_regsitry(),
            "tag": self.generate_matrix_tags(),
        }

    def generate_matrix_tags(self):
        for cuda in SUPPORTED_CUDA_LIST:
            for python in SUPPORTED_PYVERSION_LIST:
                for golang in SUPPORTED_GOLANG_LIST:
                    for distro in SUPPORTED_DISTRO_LIST:
                        self.matrix.append(ImageTag(cuda, python=python, golang=golang, distro=distro)
                                           )

    def generate_dockerfile(self):
        if not os.path.exists(TEMP_GEN_DIR):
            os.makedirs(TEMP_GEN_DIR)
        for tag in self.matrix:
            dockerfile_dir = os.path.join(TEMP_GEN_DIR, tag.distro,
                                          "cuda-" + tag.cuda)  # for now , we fixed python version and golang
            st = self.render(tag)
            if not os.path.exists(dockerfile_dir):
                os.makedirs(dockerfile_dir)
            with open(os.path.join(dockerfile_dir, Dockerfile), 'w') as dockerfile:
                dockerfile.write(st)
                dockerfile.close()
                log.info("write %s success" % os.path.abspath(os.path.join(dockerfile_dir, Dockerfile)))

    def render(self, tag: ImageTag):
        s = self.template.render(use_github=self.use_github, vars={
            "registry": self.get_regsitry(),
            "tag": str(tag)
        })
        return s

    def get_regsitry(self):
        if self.use_github:
            return ECR_REPO
        return INNER_REPO

    def set_output_path(self, target):
        self.output_path = pathlib.Path(
            f"{self.dist_base_path}/{target.replace('.', '')}"
        )
        if not self.parent.shipit_uuid and self.output_path.exists:
            log.info(f"Removing {self.output_path}")
            rm["-rf", self.output_path]()
        log.debug(f"self.output_path: '{self.output_path}' target: '{target}'")
        log.debug(f"Creating {self.output_path}")
        self.output_path.mkdir(parents=True, exist_ok=False)

    def _load_template(self):
        tpl = "./docker/templates/aiges-gpu/Dockerfile.j2"
        if not os.path.exists(tpl):
            raise FileNotFoundError("not found %s" % tpl)
        log.info("load success j2 file.")
        self.template = self.template_env.from_string(open(tpl, "r").read())

    def targeted(self):
        self._load_template()
        self.generate_matrix_tags()
        self.generate_dockerfile()

    def main(self):
        self.targeted()
        log.info("Done")


if __name__ == "__main__":
    Manager.run()
